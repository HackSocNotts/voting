package common

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type VoterRegistration struct {
	BallotID  string
	StudentID int
}

func Connect() (*mongo.Client, error) {
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
