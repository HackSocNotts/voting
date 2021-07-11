package main

import (
	"context"
	"encoding/json"
	"log"
	"math"
	"net/http"

	"hacksocnotts.co.uk/voting/common"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Client

type Results struct {
	Winners    map[string]string `json:"winners"`
	NumBallots int               `json:"num_ballots"`
	NumVotes   int               `json:"num_votes"`
}

func main() {
	var err error

	db, err = common.Connect()
	if err != nil {
		log.Fatal("could not connect to the database.", err)
	}

	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	r.PathPrefix("/results").HandlerFunc(handleResults)
	r.Path("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})

	log.Println("starting admin control panel server on :10002")
	log.Fatal(http.ListenAndServe(":10002", r))
}

func handleResults(w http.ResponseWriter, r *http.Request) {
	ballots, err := getBallots()
	if err != nil {
		log.Println("there was an error receiving the ballots.", err)
		common.Error(w, http.StatusInternalServerError, "There was a database error, please try again.")
		return
	}

	candidates, err := getCandidates()
	if err != nil {
		log.Println("there was an error receiving the candidates.", err)
		common.Error(w, http.StatusInternalServerError, "There was a database error, please try again.")
		return
	}

	res := &Results{
		Winners:    make(map[string]string),
		NumBallots: len(ballots),
		NumVotes:   0,
	}

	for _, b := range ballots {
		if b.Votes != nil {
			res.NumVotes++
		}
	}

	for i, pos := range candidates {
		winner := calculateWinner(i, ballots)
		res.Winners[pos.Role] = pos.Candidates[winner]
	}

	resp, err := json.Marshal(res)
	if err != nil {
		common.Error(w, http.StatusInternalServerError, "There was an unexpected error, please try again.")
		return
	}

	w.Write(resp)
}

func calculateWinner(pos int, ballots []common.Ballot) int {
	matrix := make([][]int, len(ballots))
	for i, b := range ballots {
		matrix[i] = (*b.Votes)[pos]
	}

	if len(ballots) == 0 {
		return -1
	}

	for len(matrix[0]) > 1 {
		// eliminate one candidate
		loser := findLoser(matrix)

		for i, row := range matrix {
			newRow := make([]int, 0, len(row)-1)
			for _, c := range row {
				if c != loser {
					newRow = append(newRow, c)
				}
			}
			matrix[i] = newRow
		}
	}

	return matrix[0][0]
}

func findLoser(m [][]int) int {
	counts := make(map[int]int)

	for _, c := range m[0] {
		counts[c] = 0
	}

	for _, b := range m {
		counts[b[0]]++
	}

	var (
		loser = -1
		min   = math.MaxInt32
	)

	for c, v := range counts {
		if loser == -1 || v < min {
			loser = c
			min = v
		}
	}

	return loser
}

// counts the amount of votes of each candidate for each position.
// if x is the return value, x[i][j] is the amount of votes of candidate #j in position #i
func countVotes(ballots []common.Ballot) [][]int {
	var (
		res [][]int
	)

	for _, b := range ballots {
		if b.Votes == nil {
			continue
		}

		if len(res) == 0 {
			res = make([][]int, len(*b.Votes))
			for i, v := range *b.Votes {
				res[i] = make([]int, len(v))
			}
		}

		for i, v := range *b.Votes {
			res[i][v[0]]++
		}
	}

	return res
}

func getBallots() ([]common.Ballot, error) {
	var (
		collection  = db.Database("Hacksoc").Collection("ballots")
		ballots     []common.Ballot
		cursor, err = collection.Find(context.TODO(), bson.D{})
	)

	if err != nil {
		return nil, err
	}

	cursor.All(context.TODO(), &ballots)

	return ballots, nil
}

func getCandidates() ([]common.Position, error) {
	var (
		collection  = db.Database("Hacksoc").Collection("candidates")
		candidates  []common.Position
		opts        = options.Find().SetSort(bson.D{{Key: "index", Value: 1}})
		cursor, err = collection.Find(context.TODO(), bson.D{}, opts)
	)

	if err != nil {
		return nil, err
	}

	cursor.All(context.TODO(), &candidates)

	return candidates, nil
}
