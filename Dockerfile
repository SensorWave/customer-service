FROM golang:1.23-alpine

# Installer git et bash
RUN apk add --no-cache git bash

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o customer-service ./cmd/main.go

EXPOSE 8080
CMD ["./customer-service"]
