package main

import (
	_ "github.com/lib/pq"

	"log"
	"./storage"
	"./server"
)

//var dbSource = os.Getenv("DATABASE_URL")
var dbSource = "user=postgres password=1234 dbname=mash sslmode=disable"

func main() {

	//port := os.Getenv("PORT")
	port := "8080"

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	storage, err := storage.Connect("postgres", dbSource)

	if err != nil {
		log.Fatal(err)
	}

	server := server.New(storage)

	err = server.Start(port)
	if err != nil {
		log.Fatal(err)
	}
	storage.DB.Close()
}
