# Variables
BINARY_NAME=quanto-eu-vou-pagar
MAIN_FILE=cmd/server/main.go

# Colors
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput -Txterm sgr0)

.PHONY: all test build run clean

all: help

## help: Imprime este help.
help:
	@echo '${GREEN}Usage:${RESET}'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo '${GREEN}Targets:${RESET}'
	@awk '/^[a-zA-Z\-0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($1, 0, index($1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  ${YELLOW}%-16s${RESET} %s\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $0 }' $(MAKEFILE_LIST)

## test: Executa os testes unitários.
test:
	@echo '${GREEN}Executando testes...${RESET}'
	go test ./... -v

## test-cover: Executa os testes unitários com cobertura.
test-cover:
	@echo '${GREEN}Executando testes com cobertura...${RESET}'
	go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out -o coverage.html
	@echo '${GREEN}Relatório de cobertura gerado: coverage.html${RESET}'

## test-all: Executa todos os testes (Docker, API, etc.).
test-all:
	@echo '${GREEN}Executando todos os testes...${RESET}'
	./scripts/test-all.sh

## build: Compila o binário do serviço.
build:
	@echo '${GREEN}Compilando...${RESET}'
	go build -o ${BINARY_NAME} ${MAIN_FILE}

## run: Executa o serviço.
run: build
	@echo '${GREEN}Executando o serviço...${RESET}'
	./${BINARY_NAME}

## clean: Remove o binário compilado e arquivos de cobertura.
clean:
	@echo '${GREEN}Limpando...${RESET}'
	rm -f ${BINARY_NAME} coverage.out coverage.html

## fmt: Formata o código Go.
fmt:
	@echo '${GREEN}Formatando código...${RESET}'
	go fmt ./...

## vet: Examina o código Go em busca de problemas.
vet:
	@echo '${GREEN}Examinando código...${RESET}'
	go vet ./...

## tidy: Atualiza go.mod e go.sum.
tidy:
	@echo '${GREEN}Atualizando dependências...${RESET}'
	go mod tidy

## docker-up: Inicia os serviços com docker-compose.
docker-up:
	@echo '${GREEN}Iniciando serviços com docker-compose...${RESET}'
	docker-compose up -d

## docker-down: Para os serviços com docker-compose.
docker-down:
	@echo '${GREEN}Parando serviços com docker-compose...${RESET}'
	docker-compose down

## docker-down-v: Para os serviços com docker-compose e remove os volumes.
docker-down-v:
	@echo '${GREEN}Parando serviços com docker-compose e removendo volumes...${RESET}'
	docker-compose down -v

## docker-logs: Mostra os logs dos serviços.
docker-logs:
	@echo '${GREEN}Mostrando logs dos serviços...${RESET}'
	docker-compose logs -f