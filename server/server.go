package server

import (
	"github.com/DiGregory/facemash/storage"
	"net/http"
	"fmt"
)

type Server struct {
	s *storage.Storage
}

func (s *Server) Start(port string) (err error) {
	http.HandleFunc("/", s.mainPageHandler)
	http.HandleFunc("/vote/", s.voteHandler)
	http.Handle("/tmp/", http.StripPrefix("/tmp/", http.FileServer(http.Dir("./../tmp")))) //girls photos

	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		return err
	}
	fmt.Println("server started")
	return

}

func New(s *storage.Storage) *Server {
	return &Server{
		s: s,
	}
}
