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
	router.HandleFunc("/sessions/{code}/join", h.JoinSession).Methods("POST")
	router.HandleFunc("/sessions/{code}", h.GetSessionByCode).Methods("GET")
	router.HandleFunc("/sessions/{code}/state", h.UpdateSessionState).Methods("PUT")
	router.HandleFunc("/sessions/{code}/leave", h.LeaveSession).Methods("POST")
	router.HandleFunc("/sessions/{code}/cards", h.CreateCardInSession).Methods("POST")
}

func (h *SessionHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	response, err := h.service.CreateSession(req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, response)
}

func (h *SessionHandler) JoinSession(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var req domain.JoinSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, err := h.service.JoinSession(params["code"], req)
	if err != nil {
		switch err {
		case service.ErrSessionNotFound:
			respondWithError(w, http.StatusNotFound, err.Error())
		case service.ErrSessionClosed:
			respondWithError(w, http.StatusForbidden, err.Error())
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}

func (h *SessionHandler) UpdateSessionState(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := r.Header.Get("User-ID") // Assumindo que o ID do usuário vem no header
	if userID == "" {
		respondWithError(w, http.StatusUnauthorized, "User ID is required")
		return
	}

	var req domain.UpdateSessionStateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err := h.service.UpdateSessionState(params["code"], userID, req)
	if err != nil {
		switch err {
		case service.ErrSessionNotFound:
			respondWithError(w, http.StatusNotFound, err.Error())
		case service.ErrUnauthorized:
			respondWithError(w, http.StatusForbidden, err.Error())
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Session state updated successfully"})
}

func (h *SessionHandler) LeaveSession(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := r.Header.Get("User-ID")
	if userID == "" {
		respondWithError(w, http.StatusUnauthorized, "User ID is required")
		return
	}

	err := h.service.LeaveSession(params["code"], userID)
	if err != nil {
		switch err {
		case service.ErrSessionNotFound:
			respondWithError(w, http.StatusNotFound, err.Error())
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Left session successfully"})
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
	userID := r.Header.Get("User-ID")
	if userID == "" {
		respondWithError(w, http.StatusUnauthorized, "User ID is required")
		return
	}

	var card domain.Card
	if err := json.NewDecoder(r.Body).Decode(&card); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	card, err := h.service.CreateCardInSession(params["code"], userID, card)
	if err != nil {
		switch err {
		case service.ErrSessionNotFound:
			respondWithError(w, http.StatusNotFound, err.Error())
		case service.ErrUnauthorized:
			respondWithError(w, http.StatusForbidden, err.Error())
		case service.ErrSessionClosed:
			respondWithError(w, http.StatusForbidden, err.Error())
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusCreated, card)
}
