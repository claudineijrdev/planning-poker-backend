.PHONY: dev build run test clean install-tools

# Variáveis
BINARY_NAME=flash-cards-api
BUILD_DIR=./build
GOPATH=$(shell go env GOPATH)
AIR=$(GOPATH)/bin/air

# Comandos principais
dev: install-air
	$(AIR)

build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/api

run: build
	./$(BUILD_DIR)/$(BINARY_NAME)

test:
	go test -v ./...

clean:
	go clean
	rm -rf $(BUILD_DIR)
	rm -rf tmp

# Instalação de ferramentas
install-tools: install-air

install-air:
	@if ! command -v $(AIR) > /dev/null; then \
		echo "Installing air..."; \
		go install github.com/cosmtrek/air@v1.49.0; \
	fi

# Ajuda
help:
	@echo "Comandos disponíveis:"
	@echo "  make dev         - Inicia o servidor com live reload usando air"
	@echo "  make build      - Compila o projeto"
	@echo "  make run        - Compila e executa o projeto"
	@echo "  make test       - Executa os testes"
	@echo "  make clean      - Remove arquivos temporários e compilados"
	@echo "  make install-tools - Instala ferramentas necessárias (air)"

# Comando padrão
.DEFAULT_GOAL := help 