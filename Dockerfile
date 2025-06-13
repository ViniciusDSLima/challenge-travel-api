# Estágio de compilação
FROM golang:1.24 AS builder

WORKDIR /app

# Copiar os arquivos go.mod e go.sum para baixar dependências
COPY go.mod go.sum ./
RUN go mod download

# Copiar o código fonte
COPY . .

# Compilar a aplicação principal
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# Estágio final - usar a imagem Go completa em vez de alpine
FROM golang:1.24

WORKDIR /app

# Instalar certificados para HTTPS
RUN apt-get update && apt-get install -y ca-certificates

# Copiar o código fonte (necessário para go run)
COPY . .

# Copiar o binário compilado do estágio de compilação
COPY --from=builder /app/main .

# Expor a porta que a aplicação usa
EXPOSE 8080

# Comando para executar a aplicação
CMD ["./main"]