package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"hacksocnotts.co.uk/voting/common"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/google/uuid"
)

var db *mongo.Client

func main() {
	var err error

	db, err = connect()
	if err != nil {
		log.Fatal("could not connect to the database.", err)
	}

	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
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
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "An unexpected error occurred")
		return
	}

	id, err := strconv.ParseInt(string(body), 10, 32)
	if err != nil {
		log.Println("error parsing given id:", string(body))
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Invalid ID format")
		return
	}

	log.Println("attempting to register user", id)
	exists, err := verify(int(id))
	if err != nil {
		log.Printf("user %d: couldn't verify user against the database. %s\n", id, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "There was a database error, please try again.")
		return
	}

	if !exists {
		log.Printf("user %d attempted to register but isn't in HackSoc\n", id)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "You don't appear to be a member of HackSoc. Are you sure you entered your ID correctly?")
		return
	}

	votedAlready, err := hasVoted(int(id))
	if err != nil {
		log.Printf("user %d: couldn't check if they have already voted. %s\n", id, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "There was a database error, please try again.")
		return
	}

	if votedAlready {
		log.Printf("user %d has already registered to vote\n", id)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "You've already registered a ballot! If this wasn't you, please contact the committee.")
		return
	}

	ballotID, err := register(int(id))
	if err != nil {
		log.Printf("user %d: couldn't register a new document in the database. %s\n", id, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "There was a database error, please try again.")
		return
	}

	log.Printf("registered user %d, ballot id %s\n", id, ballotID)
}

func connect() (*mongo.Client, error) {
	var (
		username = os.Getenv("MONGO_USER")
		password = os.Getenv("MONGO_PASS")
	)

	opts := options.Client().ApplyURI(fmt.Sprintf("mongodb+srv://%s:%s@cluster0.q5uor.mongodb.net/hacksoc?retryWrites=true&w=majority", username, password))

	log.Println("connecting to mongodb database...")
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, err
	}
	log.Println("connected to database")

	err = client.Ping(context.TODO(), nil)

	return client, err
}

func verify(id int) (bool, error) {
	var (
		filter     = bson.D{{"ID", id}}
		collection = db.Database("Hacksoc").Collection("members")
		n, err     = collection.CountDocuments(context.TODO(), filter)
	)

	return n > 0, err
}

func hasVoted(id int) (bool, error) {
	var (
		filter     = bson.D{{"studentid", id}}
		collection = db.Database("Hacksoc").Collection("voting_reg")
		n, err     = collection.CountDocuments(context.TODO(), filter)
	)

	return n > 0, err
}

func register(id int) (string, error) {
	var (
		ballotID = uuid.New()
		reg      = common.VoterRegistration{
			BallotID:  ballotID.String(),
			StudentID: id,
		}
		collection = db.Database("Hacksoc").Collection("voting_reg")
	)

	_, err := collection.InsertOne(context.TODO(), reg)
	return ballotID.String(), err
}
