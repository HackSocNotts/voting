package common

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MemberVoted struct {
	StudentID int
}

type Ballot struct {
	Votes *map[int]([]int) `json:"votes"`
}

type Position struct {
	Index      int      `json:"index"`
	Role       string   `json:"role"`
	Candidates []string `json:"candidates"`
}

func Connect() (*mongo.Client, error) {
	var (
		username = os.Getenv("MONGO_USER")
		password = os.Getenv("MONGO_PASS")
		host = os.Getenv("MONGO_HOST")
	)

	opts := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:27017/hacksoc?retryWrites=true&w=majority", username, password, host))

	log.Println("connecting to mongodb database...")
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, err
	}
	log.Println("connected to database")

	err = client.Ping(context.TODO(), nil)

	return client, err
}

func Error(w http.ResponseWriter, status int, msg string, args ...interface{}) {
	w.WriteHeader(status)
	fmt.Fprintf(w, msg, args...)
}
