package handler

import (
	"net/http"

	"flash-cards/backend/internal/service"
	"flash-cards/backend/internal/websocket"

	"github.com/gorilla/mux"
)

// WebsocketHandler gerencia as conexões WebSocket
type WebsocketHandler struct {
	websocketService *service.WebsocketService
}

// NewWebsocketHandler cria uma nova instância do handler de WebSocket
func NewWebsocketHandler(websocketService *service.WebsocketService) *WebsocketHandler {
	return &WebsocketHandler{
		websocketService: websocketService,
	}
}

// RegisterRoutes registra as rotas do WebSocket
func (h *WebsocketHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/ws/{sessionCode}", h.HandleWebSocket).Methods("GET")
}

// HandleWebSocket gerencia a conexão WebSocket para uma sessão específica
func (h *WebsocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionCode := vars["sessionCode"]
	
	hub := h.websocketService.GetHub(sessionCode)
	websocket.ServeWs(hub, w, r)
} 