package service

import (
	"sync"

	"flash-cards/backend/internal/domain"
	"flash-cards/backend/internal/websocket"
)

// WebsocketService gerencia as conexões WebSocket e broadcasts
type WebsocketService struct {
	hubs map[string]*websocket.Hub // Mapeia códigos de sessão para hubs
	mutex sync.Mutex
}

// NewWebsocketService cria uma nova instância do serviço de WebSocket
func NewWebsocketService() *WebsocketService {
	return &WebsocketService{
		hubs: make(map[string]*websocket.Hub),
	}
}

// GetHub retorna o hub para uma sessão específica, criando um novo se não existir
func (s *WebsocketService) GetHub(sessionCode string) *websocket.Hub {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if hub, exists := s.hubs[sessionCode]; exists {
		return hub
	}

	hub := websocket.NewHub()
	s.hubs[sessionCode] = hub
	go hub.Run()
	return hub
}

// RemoveHub remove um hub quando a sessão é fechada
func (s *WebsocketService) RemoveHub(sessionCode string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if hub, exists := s.hubs[sessionCode]; exists {
		// Fechar todos os clientes
		for client := range hub.Clients() {
			hub.Unregister(client)
		}
		delete(s.hubs, sessionCode)
	}
}

// BroadcastSession envia uma atualização da sessão para todos os clientes conectados
func (s *WebsocketService) BroadcastSession(session domain.Session) {
	hub := s.GetHub(session.Code)
	hub.Broadcast(session)
}

// BroadcastCard envia uma atualização de card para todos os clientes conectados à sessão
func (s *WebsocketService) BroadcastCard(sessionCode string, card domain.Card) {
	hub := s.GetHub(sessionCode)
	hub.Broadcast(card)
}

// BroadcastUserUpdate envia uma atualização de usuário para todos os clientes conectados à sessão
func (s *WebsocketService) BroadcastUserUpdate(sessionCode string, user domain.User, action string) {
	message := map[string]interface{}{
		"type": "user_update",
		"action": action,
		"user": user,
	}
	
	hub := s.GetHub(sessionCode)
	hub.Broadcast(message)
} 