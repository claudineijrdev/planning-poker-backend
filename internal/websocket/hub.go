package websocket

import (
	"encoding/json"
	"log"
	"sync"
)

// Hub mantém o conjunto de conexões WebSocket ativas
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mutex      sync.Mutex
}

// Client representa uma conexão WebSocket
type Client struct {
	hub  *Hub
	send chan []byte
}

// NewHub cria uma nova instância do Hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run inicia o hub e gerencia as conexões
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client] = true
			h.mutex.Unlock()

		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mutex.Unlock()

		case message := <-h.broadcast:
			h.mutex.Lock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mutex.Unlock()
		}
	}
}

// Broadcast envia uma mensagem para todos os clientes conectados
func (h *Hub) Broadcast(message interface{}) {
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("Erro ao serializar mensagem: %v", err)
		return
	}
	h.broadcast <- jsonMessage
}

// Clients retorna todos os clientes conectados
func (h *Hub) Clients() map[*Client]bool {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	
	// Criar uma cópia do mapa para evitar problemas de concorrência
	clientsCopy := make(map[*Client]bool)
	for client := range h.clients {
		clientsCopy[client] = true
	}
	return clientsCopy
}

// Unregister remove um cliente do hub
func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
} 