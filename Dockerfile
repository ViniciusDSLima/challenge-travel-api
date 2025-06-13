FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

FROM golang:1.24

WORKDIR /app

RUN apt-get update && apt-get install -y ca-certificates

COPY . .

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]