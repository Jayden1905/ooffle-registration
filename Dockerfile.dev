FROM golang:1.23.1-alpine AS base

WORKDIR /app

# Install Air for hot reloading
RUN go install github.com/air-verse/air@latest
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

RUN apk add --no-cache bash curl

# Expose the application port
EXPOSE 8080

# Use Air for hot reloading
CMD ["bash", "-c", "goose -dir ./cmd/sql/migrations mysql \"${DB_USER}:${DB_PASSWD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}\" up && air -c .air.toml"]
