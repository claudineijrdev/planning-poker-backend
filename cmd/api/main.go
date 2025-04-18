package main

import (
	"log"
	"net/http"

	"flash-cards/backend/internal/handler"
	"flash-cards/backend/internal/repository"
	"flash-cards/backend/internal/service"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// Inicialização dos repositórios
	cardRepo := repository.NewCardRepository()
	sessionRepo := repository.NewSessionRepository()

	// Inicialização dos serviços
	cardService := service.NewCardService(cardRepo)
	sessionService := service.NewSessionService(sessionRepo, cardRepo)
	websocketService := service.NewWebsocketService()

	// Inicialização dos handlers
	cardHandler := handler.NewCardHandler(cardService, websocketService)
	sessionHandler := handler.NewSessionHandler(sessionService, websocketService)
	websocketHandler := handler.NewWebsocketHandler(websocketService)

	// Configuração do router
	router := mux.NewRouter()

	// Registro das rotas
	cardHandler.RegisterRoutes(router)
	sessionHandler.RegisterRoutes(router)
	websocketHandler.RegisterRoutes(router)

	// Configuração do CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:5173", "http://127.0.0.1:5500"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "User-ID", "Authorization", "Origin"},
		AllowCredentials: true,
		Debug:            true,
	})

	// Inicialização do servidor
	handler := c.Handler(router)
	log.Println("Server starting on port 3001...")
	log.Fatal(http.ListenAndServe(":3001", handler))
}
