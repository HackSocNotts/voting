package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"hacksocnotts.co.uk/voting/common"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var db *mongo.Client

func main() {
	var err error

	db, err = common.Connect()
	if err != nil {
		log.Fatal("could not connect to the database.", err)
	}

	r := mux.NewRouter().PathPrefix("/register").Subrouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/register/static/", http.FileServer(http.Dir("./static/"))))
	r.Path("/register/").Methods("POST").HandlerFunc(registerHandler)
	r.Path("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})

	log.Println("starting registration server on :10000")
	log.Fatal(http.ListenAndServe(":10000", r))
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("error reading /register req. body:", err)
		common.Error(w, http.StatusBadRequest, "An unexpected error occurred")
		return
	}

	id, err := strconv.ParseInt(string(body), 10, 32)
	if err != nil {
		log.Println("error parsing given id:", string(body))
		common.Error(w, http.StatusBadRequest, "Invalid ID format. This shouldn't happen, please contact the committee.")
		return
	}

	log.Println("attempting to register user", id)
	exists, err := verify(int(id))
	if err != nil {
		log.Printf("user %d: couldn't verify user against the database. %s\n", id, err.Error())
		common.Error(w, http.StatusInternalServerError, "There was a database error, please try again.")
		return
	}

	if !exists {
		log.Printf("user %d attempted to register but isn't in HackSoc\n", id)
		common.Error(w, http.StatusBadRequest, "You don't appear to be a member of HackSoc. Are you sure you entered your ID correctly?")
		return
	}

	votedAlready, err := hasVoted(int(id))
	if err != nil {
		log.Printf("user %d: couldn't check if they have already voted. %s\n", id, err.Error())
		common.Error(w, http.StatusInternalServerError, "There was a database error, please try again.")
		return
	}

	if votedAlready {
		log.Printf("user %d has already registered to vote\n", id)
		common.Error(w, http.StatusBadRequest, "You've already registered a ballot! If this wasn't you, please contact the committee.")
		return
	}

	ballotID, err := register(int(id))
	if err != nil {
		log.Printf("user %d: couldn't register a new document in the database. %s\n", id, err.Error())
		common.Error(w, http.StatusInternalServerError, "There was a database error, please try again.")
		return
	}

	log.Printf("registered user %d\n", id)

	fmt.Fprint(w, ballotID)
}

func verify(id int) (bool, error) {
	var (
		filter     = bson.D{{Key: "ID", Value: id}}
		collection = db.Database("Hacksoc").Collection("members")
		n, err     = collection.CountDocuments(context.TODO(), filter)
	)

	return n > 0, err
}

func hasVoted(id int) (bool, error) {
	var (
		filter     = bson.D{{Key: "studentid", Value: id}}
		collection = db.Database("Hacksoc").Collection("members_voted")
		n, err     = collection.CountDocuments(context.TODO(), filter)
	)

	return n > 0, err
}

func register(id int) (string, error) {
	var (
		membersVoted = db.Database("Hacksoc").Collection("members_voted")
		ballots      = db.Database("Hacksoc").Collection("ballots")
	)

	_, err := membersVoted.InsertOne(context.TODO(), common.MemberVoted{StudentID: id})
	if err != nil {
		return "", err
	}

	res, err := ballots.InsertOne(context.TODO(), common.Ballot{})
	if err != nil {
		return "", err
	}

	switch ballotID := res.InsertedID.(type) {
	case primitive.ObjectID:
		return ballotID.Hex(), nil
	default:
		return "", errors.New("unexpected error. ballot ID is not an ObjectID")
	}
}
