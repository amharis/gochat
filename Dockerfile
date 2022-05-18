# syntax=docker/dockerfile:1
FROM golang:1.18-bullseye

WORKDIR /app
COPY ../* .
CMD go run *.go -chat-backend=$CHATROOM_LISTENER -redis-connection-string=$REDIS_HOST:6379

