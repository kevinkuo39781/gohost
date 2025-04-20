# --- Build Stage ---
FROM golang:1.24.2 AS builder

WORKDIR /app
COPY . .

RUN go build -o gohost ./cmd/gohost

# --- Runtime Stage ---
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/gohost .
COPY public ./public

RUN mkdir -p ./data/sites

EXPOSE 8080
CMD ["./gohost"]