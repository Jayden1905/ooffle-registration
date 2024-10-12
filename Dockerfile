FROM golang:1.23.1-alpine AS builder  

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/bin/event-registration-software ./cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN apk add --no-cache bash curl

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /app/bin/event-registration-software .
COPY wait-for-it.sh /app/wait-for-it.sh

RUN chown -R appuser:appgroup /app
USER appuser

EXPOSE 8080

CMD ["/app/wait-for-it.sh", "db:3306", "--", "./event-registration-software"]
