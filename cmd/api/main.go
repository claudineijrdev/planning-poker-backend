package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"flash-cards/backend/internal/handler"
	"flash-cards/backend/internal/repository"
	"flash-cards/backend/internal/service"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// Inicialização do gerador de números aleatórios
	rand.Seed(time.Now().UnixNano())

	// Inicialização dos repositórios
	cardRepo := repository.NewCardRepository()
	sessionRepo := repository.NewSessionRepository()

	// Inicialização dos serviços
	cardService := service.NewCardService(cardRepo)
	sessionService := service.NewSessionService(sessionRepo, cardRepo)

	// Inicialização dos handlers
	cardHandler := handler.NewCardHandler(cardService)
	sessionHandler := handler.NewSessionHandler(sessionService)

	// Configuração do router
	router := mux.NewRouter()

	// Registro das rotas
	cardHandler.RegisterRoutes(router)
	sessionHandler.RegisterRoutes(router)

	// Configuração do CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
	})

	// Inicialização do servidor
	handler := c.Handler(router)
	log.Println("Server starting on port 3001...")
	log.Fatal(http.ListenAndServe(":3001", handler))
}
