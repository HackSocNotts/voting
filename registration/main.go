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

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	log.Println("registering user", id)
	exists, err := verify(int(id))
	if err != nil {
		log.Printf("user %d: couldn't verify user against the database\n", id)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "There was a database error, please try again.")
		return
	}

	if !exists {
		log.Printf("user %d attempted to register but isn't in HackSoc\n", id)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "You don't appear to be a member of HackSoc. Are you sure you entered your ID correctly?")
		return
	}
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

}
