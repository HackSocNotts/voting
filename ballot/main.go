package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"hacksocnotts.co.uk/voting/common"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Client

type Position struct {
	Index      int      `json:"index"`
	Role       string   `json:"role"`
	Candidates []string `json:"candidates"`
}

type BallotRequest struct {
	Ballot common.Ballot `json:"ballot"`
	ID     string        `json:"id"`
}

func main() {
	var err error

	db, err = common.Connect()
	if err != nil {
		log.Fatal("could not connect to the database.", err)
	}

	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	r.Path("/candidates/").HandlerFunc(handleCandidates)
	r.Path("/submit/").Methods("POST").HandlerFunc(handleSubmit)
	r.Path("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})

	log.Println("starting ballot server on :10001")
	log.Fatal(http.ListenAndServe(":10001", r))
}

func handleCandidates(w http.ResponseWriter, r *http.Request) {
	var (
		collection  = db.Database("Hacksoc").Collection("candidates")
		opts        = options.Find().SetSort(bson.D{{Key: "index", Value: 1}})
		cursor, err = collection.Find(context.TODO(), bson.D{}, opts)
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

func handleSubmit(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("error reading /submit req. body:", err)
		common.Error(w, http.StatusBadRequest, "An unexpected error occurred")
		return
	}

	var ballot BallotRequest

	json.Unmarshal(body, &ballot)

	collection := db.Database("Hacksoc").Collection("ballots")
	id, err := primitive.ObjectIDFromHex(ballot.ID)
	if err != nil {
		log.Printf("invalid _id %s (%s)\n", id, err.Error())
		common.Error(w, http.StatusBadRequest, "Your ballot ID seems to be invalid.")
		return
	}

	var existing common.Ballot

	err = collection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: id}}).Decode(&existing)
	if err != nil {
		log.Println("error checking for an existing ballot entry, id", ballot.ID, "-", err)
		common.Error(w, http.StatusInternalServerError, "There was a database error, please try again.")
		return
	}

	if existing.Votes != nil {
		log.Printf("repeated ballot submission for id %s\n", ballot.ID)
		common.Error(w, http.StatusBadRequest, "It looks like you've already submitted this ballot. If this wasn't you, please contact the committee.")
		return
	}

	log.Printf("recording ballot from ballot id %s\n", ballot.ID)

	_, err = collection.ReplaceOne(context.TODO(), bson.D{{Key: "_id", Value: id}}, ballot.Ballot)
	if err != nil {
		log.Printf("error recording a ballot in the database for id %s (%s)\n", ballot.ID, err)
		common.Error(w, http.StatusInternalServerError, "There was an error recording your ballot in the database, please try again.")
		return
	}
}
