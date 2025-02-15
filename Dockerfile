FROM golang:1.23.6-alpine

WORKDIR /app

# Install PostgreSQL client
RUN apk add --no-cache postgresql-client
COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY .env.example .env


RUN go build -o main src/telegram-bot/telegram.go
RUN mkdir -p /data

COPY data/db.sql /data/

ENV TG_BOT_TOKEN=""
ENV DATABASE_URL=""

CMD ["./main"]
