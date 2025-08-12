# Etapa de build
FROM golang:1.22.0-alpine AS builder

# Instalar dependências necessárias para o build (se houver)
# RUN apk add --no-cache git

# Definir o diretório de trabalho dentro do container
WORKDIR /app

# Copiar go mod e sum para o diretório de trabalho
COPY go.mod go.sum ./

# Baixar as dependências
RUN go mod download

# Copiar o código fonte para o diretório de trabalho
COPY . .

# Compilar o binário
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/server/main.go

# Etapa final
FROM alpine:latest

# Instalar ca-certificates para chamadas HTTPS (se necessário)
RUN apk --no-cache add ca-certificates

# Criar um diretório para a aplicação
WORKDIR /root/

# Copiar o binário do estágio de build
COPY --from=builder /app/main .

# Expor a porta que a aplicação irá usar
EXPOSE 8080

# Comando para executar a aplicação
CMD ["./main"]