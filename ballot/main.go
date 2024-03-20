package main

import (
	"context"
	"encoding/json"
	"fmt"
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

var (
	db         *mongo.Client
	candidates []common.Position
)

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

	if err = loadCandidates(); err != nil {
		log.Fatal("could not fetch candidate list.", err)
	}

	r := mux.NewRouter()
	r.PathPrefix("/ballot/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	r.Path("/ballot/candidates/").HandlerFunc(handleCandidates)
	r.Path("/ballot/active/{id}/").HandlerFunc(handleActive)
	r.Path("/ballot/submit/").Methods("POST").HandlerFunc(handleSubmit)
	r.Path("/ballot").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})
	r.Path("/ballot/thank-you").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/thankyou.html")
	})

	log.Println("starting ballot server on :10001")
	log.Fatal(http.ListenAndServe(":10001", r))
}

func handleCandidates(w http.ResponseWriter, r *http.Request) {
	res, err := json.Marshal(candidates)
	if err != nil {
		common.Error(w, http.StatusInternalServerError, "There was an unexpected error, please try again.")
		return
	}

	w.Write(res)
}

func handleActive(w http.ResponseWriter, r *http.Request) {
	var (
		id       = mux.Vars(r)["id"]
		oid, err = primitive.ObjectIDFromHex(id)
	)

	if err != nil {
		log.Println("invalid id", id)
		common.Error(w, http.StatusBadRequest, "Your ballot ID (%s) doesn't seem to be in the correct format.", id)
		return
	}

	collection := db.Database("Hacksoc").Collection("ballots")
	res := collection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: oid}})

	if res.Err() == mongo.ErrNoDocuments {
		fmt.Fprintf(w, "false")
		return
	}

	var ballot common.Ballot

	if err = res.Decode(&ballot); err != nil {
		log.Println("couldn't check for existence of ballot", id)
		common.Error(w, http.StatusInternalServerError, "There was a database error, please try again.")
		return
	}

	fmt.Fprintf(w, "%t", ballot.Votes == nil)
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

	if !allCandidatesRanked(ballot.Ballot) {
		log.Println("incomplete ballot submitted")
		common.Error(w, http.StatusBadRequest, "You must rank every candidate to vote.")
		return
	}

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

func loadCandidates() error {
	var (
		collection  = db.Database("Hacksoc").Collection("candidates")
		opts        = options.Find().SetSort(bson.D{{Key: "index", Value: 1}})
		cursor, err = collection.Find(context.TODO(), bson.D{}, opts)
	)

	if err != nil {
		return err
	}

	cursor.All(context.TODO(), &candidates)

	return nil
}

func allCandidatesRanked(ballot common.Ballot) bool {
	for i, pos := range candidates {
		ranking, ok := (*ballot.Votes)[i]
		if !ok {
			return false
		}

		if len(pos.Candidates) != len(ranking) {
			return false
		}
	}

	return true
}
