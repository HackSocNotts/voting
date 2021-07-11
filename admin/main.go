package main

import (
	"log"
	"net/http"

	"hacksocnotts.co.uk/voting/common"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

var db *mongo.Client

func main() {
	var err error

	db, err = common.Connect()
	if err != nil {
		log.Fatal("could not connect to the database.", err)
	}

	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	r.Path("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})

	log.Println("starting admin control panel server on :10002")
	log.Fatal(http.ListenAndServe(":10002", r))
}
