package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
)

type redisConnector struct {
	client         *redis.Client
	pushToRedis    chan Message
	pushToChatroom chan Message
	hub            *Hub
}

/**
type Message struct {
    Channel      string
    Pattern      string
    Payload      string
    PayloadSlice []string
}
*/
func (rc *redisConnector) run() {
	subscriber := rc.client.Subscribe(ctx, CHATROOM)

	for {
		msg, err := subscriber.ReceiveMessage(ctx)
		fmt.Println("Received Redis message", msg)
		if err != nil {
			fmt.Println("Error in message reading from Redis", err)
			panic(err)
		}
		m := Message{}
		if err := json.Unmarshal([]byte(msg.Payload), &m); err != nil {
			panic(err)
		}
		fmt.Println("Unmarshalled message", m)
		rc.hub.broadcast <- m
		// ...
	}
}
