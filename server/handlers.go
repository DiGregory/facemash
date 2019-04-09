package server

import (
	"net/http"
	"log"
	"fmt"
	"strings"
	"html/template"
)

var GirlsPageTmpl = template.Must(template.ParseFiles(
	"templates/Girls.html",
	"templates/RatingPage.html",
	"templates/VotingPage.html",
))

func (s *Server) mainPageHandler(resp http.ResponseWriter, req *http.Request) {

	if req.Method == http.MethodGet {

		var err error
		data := map[string]interface{}{}
		data["FirstGirl"], data["SecondGirl"], err = s.s.GetGirlsFromDb()
		if err != nil {
			log.Fatal("can`t get girls from db", err)

		}
		data["TopGirls"], err = s.s.GetTop()
		if err != nil {
			log.Fatal("can`t get girls from db", err)

		}

		err = GirlsPageTmpl.ExecuteTemplate(resp, "girls", data)
		if err != nil {

			http.Error(resp, err.Error(), http.StatusInternalServerError)
			return

		}
	}
}

func (s *Server) voteHandler(resp http.ResponseWriter, req *http.Request) {

	if req.Method == http.MethodPost {

		DoubleId := req.URL.Path[len("/vote/"):]
		WinId := strings.Split(DoubleId, "_")[0]
		LooseId := strings.Split(DoubleId, "_")[1]
		fmt.Printf("WinId: %v; LooseId: %v;\r\n", WinId, LooseId)
		s.s.UpdateGirls(WinId, LooseId)

		http.Redirect(resp, req, "/", http.StatusSeeOther)
		return
	}
	resp.WriteHeader(http.StatusMethodNotAllowed)
	//	http.Error(resp, err.Error(), http.StatusInternalServerError)
	return

}
