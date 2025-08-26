# Estágio de build
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o /main .

# Estágio final
FROM alpine:latest
WORKDIR /
COPY --from=builder /main /main
CMD ["/main"]