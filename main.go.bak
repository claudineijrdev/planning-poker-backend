package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Result struct {
	Average      float64     `json:"average"`
	Distribution map[int]int `json:"distribution"`
}

var nextCardID = 0

type Card struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Votes       []int  `json:"votes"`
	Result      Result `json:"result"`
	Closed      bool   `json:"closed"`
}

func NewCard() Card {
	card := Card{
		ID: strconv.Itoa(nextCardID),
	}
	nextCardID++
	return card
}

type Vote struct {
	Score int `json:"score"`
}

var cards = []Card{}

func getCards(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cards)
}

func vote(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var vote Vote
	_ = json.NewDecoder(r.Body).Decode(&vote)

	for i := range cards {
		if cards[i].ID == params["id"] {
			cards[i].Votes = append(cards[i].Votes, vote.Score)

			// Calcular média
			sum := 0
			for _, v := range cards[i].Votes {
				sum += v
			}
			cards[i].Result.Average = float64(sum) / float64(len(cards[i].Votes))
			cards[i].Result.Distribution = make(map[int]int)
			for _, v := range cards[i].Votes {
				cards[i].Result.Distribution[v]++
			}
			break
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cards)
}

func resetAllVotings(w http.ResponseWriter, r *http.Request) {
	for i := range cards {
		cards[i].Votes = []int{}
		cards[i].Result.Average = 0
		cards[i].Result.Distribution = make(map[int]int)
		cards[i].Closed = false
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cards)
}

func closeVoting(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	for i := range cards {
		if cards[i].ID == params["id"] {
			cards[i].Closed = true
			break
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cards)
}

func createCard(w http.ResponseWriter, r *http.Request) {
	card := NewCard()
	_ = json.NewDecoder(r.Body).Decode(&card)
	cards = append(cards, card)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cards)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/cards", getCards).Methods("GET")
	router.HandleFunc("/cards", createCard).Methods("POST")
	router.HandleFunc("/cards/{id}/close", closeVoting).Methods("POST")
	router.HandleFunc("/cards/reset-all", resetAllVotings).Methods("POST")
	router.HandleFunc("/vote/{id}", vote).Methods("POST")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
	})

	handler := c.Handler(router)
	log.Fatal(http.ListenAndServe(":8080", handler))
}
