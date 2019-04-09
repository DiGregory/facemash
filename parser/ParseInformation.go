package parser

import (
	"github.com/himidori/golang-vk-api"
	"fmt"
	"os"
	"strconv"

	"io/ioutil"
	"strings"
	"time"
	"database/sql"

	_ "github.com/lib/pq"
)

var User = " "
var Password = " "
var dbSource = "user=postgres password=1234 dbname=mash sslmode=disable"

func main() {

	client, err := vkapi.NewVKClient(vkapi.DeviceAndroid, User, Password)

	if err != nil {
		fmt.Println("error was happened with authentication: ", err)
	}
	//файл со всеми ид, откуда будем брать инфу и фотку
	allids, err := ioutil.ReadFile("GoodId.txt")

	if err != nil {
		fmt.Println(err)
	}

	girlsIds := strings.Split(string(allids), "\r\n")
	girlsIdsInt := make([]int, 0)
	for _, v := range girlsIds {
		id, _ := strconv.Atoi(v)
		girlsIdsInt = append(girlsIdsInt, id)
	}

	time.Sleep(1000 * time.Millisecond)

	users, err := client.UsersGet(girlsIdsInt)
	if err != nil {
		fmt.Println(err)
	}

	//фотки
	//for i,v:=range users {
	//	if users[i].Photo_max_orig == "https://vk.com/images/camera_400.png?ava=1" {
	//		continue
	//	}
	//
	//	fmt.Println(users[i].Photo_max_orig)
	//	PhotoUrl := users[i].Photo_max_orig
	//
	//	time.Sleep(1000 * time.Millisecond)
	//
	//	response, err := http.Get(PhotoUrl)
	//	if err!=recover(){
	//		fmt.Println(err)
	//		time.Sleep(1000 * time.Millisecond)
	//	}
	//
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	defer response.Body.Close()
	//
	//	//открыть картинку для записи
	//	ImgFile, err := os.Create("tmp\\" +
	//		"" +
	//		strconv.Itoa(v.UID) +
	//		".jpg")
	//
	//	if err != nil {
	//		log.Fatal(err)
	//
	//	}
	//
	//	_, err = io.Copy(ImgFile, response.Body)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	ImgFile.Close()
	//}

	db, err := sql.Open("postgres", dbSource)
	if err != nil {

		fmt.Println([]byte(err.Error()))
	}

	defer db.Close()

	//запись в базу ид и имен
	for i, v := range users {

		//проверка, есть ли аватарка в базе. если нет, то игнорим
		checkPhoto, err := os.Open("tmp\\" +

			strconv.Itoa(v.UID) +
			".jpg")
		defer checkPhoto.Close()

		if err != nil {
			fmt.Println(err)
			continue
		}
		wins := 0
		losses := 0
		rating := 1200

		_, err = db.Exec("INSERT INTO Girls VALUES($1, $2, $3, $4, $5, $6, $7)", i, users[i].UID, rating, wins, losses, users[i].FirstName, users[i].LastName)

		if err != nil {

			fmt.Println(err)
		}

	}
}
