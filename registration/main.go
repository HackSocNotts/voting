package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func main() {
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
		return
	}

	id, err := strconv.ParseInt(string(body), 10, 32)
	if err != nil {
		log.Println("error parsing given id:", string(body))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Println("registering user", id)
}
