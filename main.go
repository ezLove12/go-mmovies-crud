package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

var movies []Movie

func main() {
	r := mux.NewRouter()
	movies = append(movies,
		Movie{
			ID:    "1",
			Isbn:  "438227",
			Title: "Movie1",
			Director: &Director{
				FirstName: "John",
				LastName:  "Smith",
			},
		},
		Movie{
			ID:    "2",
			Isbn:  "438229",
			Title: "Movie2",
			Director: &Director{
				FirstName: "Jane",
				LastName:  "Smith",
			},
		})
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Println("Starting server on port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(movies)
	if err != nil {
		return
	}
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	isDeleted := false
	for index, movie := range movies {
		if movie.ID == params["id"] {
			isDeleted = true
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
	if !isDeleted {
		http.Error(w, "Movie not found", http.StatusNotFound)
		return
	}
	err := json.NewEncoder(w).Encode(movies)
	if err != nil {
		return
	}
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, movie := range movies {
		if movie.ID == params["id"] {
			err := json.NewEncoder(w).Encode(movie)
			if err != nil {
				return
			}
			return
		}
	}
	http.Error(w, "Movie not found", http.StatusNotFound)
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)
	movie.ID = strconv.Itoa(rand.Intn(10000000))
	movies = append(movies, movie)
	err := json.NewEncoder(w).Encode(movie)
	if err != nil {
		return
	}
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var movieDto Movie
	_ = json.NewDecoder(r.Body).Decode(&movieDto)
	for i, _ := range movies {
		//
		if movies[i].ID == params["id"] {
			if movieDto.Isbn != "" {
				movies[i].Isbn = movieDto.Isbn
			}

			if movieDto.Title != "" {
				movies[i].Title = movieDto.Title
			}

			if movieDto.Director != nil {
				movies[i].Director = movieDto.Director
			}
			if err := json.NewEncoder(w).Encode(movies[i]); err != nil {
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
				return
			}
			return
		}
	}
	http.Error(w, "Movie not found", http.StatusNotFound)
}
