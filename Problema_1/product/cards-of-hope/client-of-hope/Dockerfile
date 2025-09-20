# syntax=docker/dockerfile:1

# Etapa de build
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o client ./cmd/app

# Etapa final
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/client ./client
RUN mkdir -p /app/data
ENV SERVER_ADDR=localhost:8080
CMD ["/app/client"]
