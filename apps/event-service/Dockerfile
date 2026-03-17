FROM golang:1.26.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal

RUN go build -o event-service ./cmd/api

FROM debian:stable-slim

WORKDIR /app

COPY --from=builder /app/event-service /app/event-service

EXPOSE 8080

CMD ["/app/event-service"]