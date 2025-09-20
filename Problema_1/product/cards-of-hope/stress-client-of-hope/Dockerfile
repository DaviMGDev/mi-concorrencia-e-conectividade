# syntax=docker/dockerfile:1

FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o stress-client main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/stress-client ./stress-client
CMD ["./stress-client"]
