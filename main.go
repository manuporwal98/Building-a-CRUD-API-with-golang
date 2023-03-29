package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux" //to bypass the go.mod issue, use go get "github.com/gorill/mux" from powershell instead of VScode
)

type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:lastname"`
}

var movies []Movie

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") //to decalre that we will work with json
	json.NewEncoder(w).Encode(movies)                  //to print out all the movies in a json format
}
func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)             //whatever comes from user in r , gets saved in params to use it later
	for index, item := range movies { //for loop over index(starting from 0) as well as directly on structures
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...) //this is a way to easily delete the value at index , obviously some other way can be used
			break
		}
	}
	json.NewEncoder(w).Encode(movies)
}
func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range movies { //blank indetifier is used because we don't need the index here
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)    //this passes the body to a variable named movie
	movie.ID = strconv.Itoa(rand.Intn(100000000)) //we are creating a random ID , this is required
	movies = append(movies, movie)
	json.NewEncoder(w).Encode(movies)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") //to update a movie , we are deleting and creating a new one
	params := mux.Vars(r)
	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			var movie Movie
			_ = json.NewDecoder(r.Body).Decode(&movie)
			movie.ID = strconv.Itoa(rand.Intn(100000000))
			movies = append(movies, movie)
			json.NewEncoder(w).Encode(movies)
			return
		}
	}
}

func main() {

	r := mux.NewRouter() //we are using the mux package to create a newrouter

	movies = append(movies, Movie{ID: "1", Isbn: "12345", Title: "One", Director: &Director{Firstname: "manu", Lastname: "Porwal"}}) ///we are not using DB here, just using slices in golang
	movies = append(movies, Movie{ID: "2", Isbn: "678910", Title: "Two", Director: &Director{Firstname: "Kapila", Lastname: "Porwal"}})

	r.HandleFunc("/movies", getMovies).Methods("GET") //this function is used to handle the route which comes after localhost:8000/ and the postman method associated with it.
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")
	fmt.Printf("Starting server at port 8000")

	log.Fatal(http.ListenAndServe(":8000", r)) //used to start the server at port 8000
}
