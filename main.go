package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/njeruthuo/comms-service/emails"
)

func main() {
	db, err := DatabaseConnect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	emailsHandler := emails.DBHandler{
		DB: db,
	}

	router := mux.NewRouter()
	router.HandleFunc("/comms/send", emailsHandler.SendHandler)
	fmt.Println("Server started on port 8001")
	log.Fatal(http.ListenAndServe(":8001", router))
}
