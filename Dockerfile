# syntax=docker/dockerfile:1
FROM golang:1.18-bullseye

ENV REDIS_CONNECTION_STRING=host.docker.internal:6379
ENV CHATROOM_LISTENER=:8080

WORKDIR /app
COPY ../* .
CMD ls
CMD go run *.go -chat-backend=$CHATROOM_LISTENER -redis-connection-string=$REDIS_CONNECTION_STRING
#EXPOSE 8081, 6379
