package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
)

type redisConnector struct {
	client *redis.Client
	hub    *Hub
}

func (rc *redisConnector) run() {
	pubSub := rc.client.Subscribe(ctx, CHATROOM)

	for {
		fmt.Println("Waiting for Redis message")
		msg, err := pubSub.ReceiveMessage(ctx)
		fmt.Println("Received Redis message", msg)
		if err != nil {
			fmt.Println("Error in message reading from Redis", err)
			panic(err)
		}
		m := Message{}
		if err := json.Unmarshal([]byte(msg.Payload), &m); err != nil {
			fmt.Println("Error in parsing message received from Redis", err)
			panic(err)
		}

		rc.hub.broadcast <- m
	}
}
