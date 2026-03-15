FROM golang:1.24.4 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o /filia-app ./cmd/main.go

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /filia-app /usr/local/bin/filia-app
COPY --from=builder /app/.env /app/.env

EXPOSE 8080

CMD ["filia-app"]
