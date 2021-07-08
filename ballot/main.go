package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"hacksocnotts.co.uk/voting/common"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var db *mongo.Client

type Position struct {
	Role       string   `json:"role"`
	Candidates []string `json:"candidates"`
}

func main() {
	var err error

	db, err = common.Connect()
	if err != nil {
		log.Fatal("could not connect to the database.", err)
	}

	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	r.PathPrefix("/candidates").HandlerFunc(handleCandidates)
	r.Path("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})

	log.Println("starting ballot server on :10001")
	log.Fatal(http.ListenAndServe(":10001", r))
}

func handleCandidates(w http.ResponseWriter, r *http.Request) {
	var (
		collection  = db.Database("Hacksoc").Collection("candidates")
		cursor, err = collection.Find(context.TODO(), bson.D{})
		results     []Position
	)

	if err != nil {
		common.Error(w, http.StatusInternalServerError, "There was a database error, please try again.")
		return
	}

	cursor.All(context.TODO(), &results)

	res, err := json.Marshal(results)
	if err != nil {
		common.Error(w, http.StatusInternalServerError, "There was an unexpected error, please try again.")
		return
	}

	w.Write(res)
}
