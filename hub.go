// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool
	// Inbound messages from the clients.
	broadcast chan Message
	// Register requests from the clients.
	register chan *Client
	// Unregister requests from clients.
	unregister chan *Client
	//redis client
	redisConnector redisConnector
	ID             string
	publish        chan Message
}

func newHub() *Hub {
	c := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	fmt.Printf("Created new client %+v \n", c)
	rc := redisConnector{
		client: c,
	}
	h := &Hub{
		broadcast:      make(chan Message, 1024),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
		clients:        make(map[*Client]bool),
		redisConnector: rc,
		ID:             "hub-" + randomString(2),
		publish:        make(chan Message, 1024),
	}
	rc.hub = h
	go rc.run()
	return h
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				log.Println("unregistering for client: ")
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.publish:
			log.Printf("Publishing message to redis %v\n", message)
			// send Message to redis
			message_bytes, err := json.Marshal(message)
			if err != nil {
				log.Printf("Error in marshalling %s message %s", err, message)
			} else {
				code := h.redisConnector.client.Publish(ctx, CHATROOM, message_bytes)
				log.Printf("Publish code: ", code.Name(), code.Err())
			}

		case message := <-h.broadcast:
			log.Printf("Broadcasting message. %v\n", message)
			// send to local clients
			for client := range h.clients {
				if true {
					select {
					case client.send <- message:
					default:
						close(client.send)
						delete(h.clients, client)
					}
				} else {
					fmt.Println("skipping sending broadcast msg to client = ", client)
				}
			}
		}
	}
}
