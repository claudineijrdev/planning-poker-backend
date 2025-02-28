package handler

import (
	"encoding/json"
	"net/http"

	"flash-cards/backend/internal/domain"
	"flash-cards/backend/internal/service"

	"github.com/gorilla/mux"
)

type CardHandler struct {
	service *service.CardService
}

func NewCardHandler(service *service.CardService) *CardHandler {
	return &CardHandler{
		service: service,
	}
}

func (h *CardHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/cards", h.GetCards).Methods("GET")
	router.HandleFunc("/cards", h.CreateCard).Methods("POST")
	router.HandleFunc("/cards/{id}/close", h.CloseVoting).Methods("POST")
	router.HandleFunc("/cards/reset-all", h.ResetAllVotings).Methods("POST")
	router.HandleFunc("/cards/{id}/vote", h.Vote).Methods("POST")
}

func (h *CardHandler) GetCards(w http.ResponseWriter, r *http.Request) {
	cards := h.service.GetAllCards()
	respondWithJSON(w, http.StatusOK, cards)
}

func (h *CardHandler) CreateCard(w http.ResponseWriter, r *http.Request) {
	var card domain.Card
	if err := json.NewDecoder(r.Body).Decode(&card); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	card = h.service.CreateCard(card)
	respondWithJSON(w, http.StatusCreated, card)
}

func (h *CardHandler) Vote(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var vote domain.Vote
	if err := json.NewDecoder(r.Body).Decode(&vote); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	card, err := h.service.AddVote(params["id"], vote)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, card)
}

func (h *CardHandler) CloseVoting(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	card, err := h.service.CloseVoting(params["id"])
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, card)
}

func (h *CardHandler) ResetAllVotings(w http.ResponseWriter, r *http.Request) {
	h.service.ResetAllVotes()
	cards := h.service.GetAllCards()
	respondWithJSON(w, http.StatusOK, cards)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
