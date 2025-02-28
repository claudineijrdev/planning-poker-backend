package handler

import (
	"encoding/json"
	"net/http"

	"flash-cards/backend/internal/domain"
	"flash-cards/backend/internal/service"

	"github.com/gorilla/mux"
)

type SessionHandler struct {
	service *service.SessionService
}

func NewSessionHandler(service *service.SessionService) *SessionHandler {
	return &SessionHandler{
		service: service,
	}
}

func (h *SessionHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/sessions", h.CreateSession).Methods("POST")
	router.HandleFunc("/sessions/join", h.JoinSession).Methods("POST")
	router.HandleFunc("/sessions/{code}", h.GetSessionByCode).Methods("GET")
	router.HandleFunc("/sessions/{sessionId}/cards", h.CreateCardInSession).Methods("POST")
}

func (h *SessionHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	response := h.service.CreateSession()
	respondWithJSON(w, http.StatusCreated, response)
}

func (h *SessionHandler) JoinSession(w http.ResponseWriter, r *http.Request) {
	var req domain.JoinSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, err := h.service.JoinSession(req)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}

func (h *SessionHandler) GetSessionByCode(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	session, err := h.service.GetSessionByCode(params["code"])
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, session)
}

func (h *SessionHandler) CreateCardInSession(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var card domain.Card
	if err := json.NewDecoder(r.Body).Decode(&card); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	card, err := h.service.CreateCardInSession(params["sessionId"], card)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, card)
}
