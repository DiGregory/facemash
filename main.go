package facemash

import (
	_ "github.com/lib/pq"

	"net/http"
	"html/template"
	"database/sql"
	"log"
	"fmt"
	"math"

	"strings"
)

var GirlsPageTmpl = template.Must(template.ParseFiles(
	"templates/Girls.html",
	"templates/RatingPage.html",
	"templates/VotingPage.html",
))

type Girl struct {
	VkId      string
	GirlID    int
	Rating    float64
	Wins      int
	Losses    int
	FirstName string
	LastName  string
}

var dbSource = "user=postgres password=1234 dbname=mash sslmode=disable"
//две рандомные записи из бд
func getGirlsFromDb() (*Girl, *Girl, error) {

	var FirstGirl = new(Girl)
	var SecondGirl = new(Girl)

	db, err := sql.Open("postgres", dbSource)

	if err != nil {
		return nil, nil, err
	}
	defer db.Close()

	FirstRecord := db.QueryRow("SELECT * FROM Girls ORDER BY RANDOM() LIMIT 1")
	err = FirstRecord.Scan(&FirstGirl.GirlID, &FirstGirl.VkId, &FirstGirl.Rating, &FirstGirl.Wins, &FirstGirl.Losses, &FirstGirl.FirstName, &FirstGirl.LastName)

	if err != nil {
		return nil, nil, err
	}

	SecondRecord := db.QueryRow("SELECT * FROM Girls ORDER BY RANDOM() LIMIT 1")
	err = SecondRecord.Scan(&SecondGirl.GirlID, &SecondGirl.VkId, &SecondGirl.Rating, &SecondGirl.Wins, &SecondGirl.Losses, &SecondGirl.FirstName, &SecondGirl.LastName)

	if err != nil {
		return nil, nil, err
	}

	if SecondGirl.VkId == FirstGirl.VkId {

		SecondRecord := db.QueryRow("SELECT * FROM Girls ORDER BY RANDOM() LIMIT 1")

		err = SecondRecord.Scan(&SecondGirl.GirlID, &SecondGirl.VkId, &SecondGirl.Rating, &SecondGirl.Wins, &SecondGirl.Losses, &SecondGirl.FirstName, &SecondGirl.LastName)

		if err != nil {
			return nil, nil, err
		}
	}

	return FirstGirl, SecondGirl, nil
}

//список лучших
func GetTop() ([]*Girl, error) {
	db, err := sql.Open("postgres", dbSource)

	if err != nil {
		return nil, err
	}
	defer db.Close()

	girls := make([]*Girl, 0)

	girlsRecords, err := db.Query("SELECT * FROM Girls ORDER BY Rating DESC LIMIT 500")
	defer girlsRecords.Close()

	for girlsRecords.Next() {
		OneGirl := new(Girl)
		err := girlsRecords.Scan(&OneGirl.GirlID, &OneGirl.VkId, &OneGirl.Rating, &OneGirl.Wins, &OneGirl.Losses, &OneGirl.FirstName, &OneGirl.LastName)
		if err != nil {
			return nil, err
		}
		girls = append(girls, OneGirl)

	}
	return girls, nil
}

func mainPageHandler(resp http.ResponseWriter, req *http.Request) {

	if req.Method == http.MethodGet {

		var err error
		data := map[string]interface{}{}
		data["FirstGirl"], data["SecondGirl"], err = getGirlsFromDb()
		if err != nil {
			log.Fatal("can`t get girls from db", err)

		}
		data["TopGirls"], err = GetTop()
		if err != nil {
			log.Fatal("can`t get girls from db", err)

		}

		err = GirlsPageTmpl.ExecuteTemplate(resp, "girls", data)
		if err != nil {
			resp.WriteHeader(http.StatusInternalServerError)
			http.Error(resp, err.Error(), http.StatusInternalServerError)
			return

		}
	}
}

//пересчет рейтинга
func UpdateGirls(WinId, LooseId string) (err error) {
	db, err := sql.Open("postgres", dbSource)

	if err != nil {
		return err
	}
	defer db.Close()

	WinGirl := new(Girl)

	WinRecord := db.QueryRow("SELECT * FROM Girls WHERE VkId=$1", WinId)

	err = WinRecord.Scan(&WinGirl.GirlID, &WinGirl.VkId, &WinGirl.Rating, &WinGirl.Wins, &WinGirl.Losses, &WinGirl.FirstName, &WinGirl.LastName)

	if err != nil {
		return err
	}

	LooseGirl := new(Girl)

	LooseRecord := db.QueryRow("SELECT * FROM Girls WHERE VkId=$1", LooseId)

	err = LooseRecord.Scan(&LooseGirl.GirlID, &LooseGirl.VkId, &LooseGirl.Rating, &LooseGirl.Wins, &LooseGirl.Losses, &LooseGirl.FirstName, &LooseGirl.LastName)

	if err != nil {
		return err
	}

	E := 1 / (1 + math.Pow(10, (WinGirl.Rating-LooseGirl.Rating)/400))
	NewWinRating := WinGirl.Rating + 10*(1-E)
	NewLooseRating := LooseGirl.Rating - 10*(1-E)

	_, err = db.Exec("UPDATE Girls SET Rating=$1, wins=wins+1 WHERE VkId=$2", NewWinRating, WinGirl.VkId)

	if err != nil {
		return err
	}

	_, err = db.Exec("UPDATE Girls SET Rating=$1, losses=losses+1 WHERE VkId=$2", NewLooseRating, LooseGirl.VkId)
	fmt.Println("girls updated")
	if err != nil {
		return err
	}

	return nil

}

func voteHandler(resp http.ResponseWriter, req *http.Request) {

	if req.Method == http.MethodPost {

		DoubleId := req.URL.Path[len("/vote/"):]
		WinId := strings.Split(DoubleId, "_")[0]
		LooseId := strings.Split(DoubleId, "_")[1]
		fmt.Printf("WinId: %v; LooseId: %v;\r\n", WinId, LooseId)
		UpdateGirls(WinId, LooseId) //пересчет рейтинга

		http.Redirect(resp, req, "/", http.StatusSeeOther)
		return
	}
	resp.WriteHeader(http.StatusMethodNotAllowed)
	//	http.Error(resp, err.Error(), http.StatusInternalServerError)
	return

}
func main() {

	http.HandleFunc("/", mainPageHandler)
	http.HandleFunc("/vote/", voteHandler)
	http.Handle("/tmp/", http.StripPrefix("/tmp/", http.FileServer(http.Dir("./tmp")))) //girls photos

	fmt.Println("server started")
	http.ListenAndServe("8080", nil)

}
