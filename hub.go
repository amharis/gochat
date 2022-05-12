// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"time"
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
			m := Message{
				Name:    "",
				Message: fmt.Sprintf("%s has joined chat", client.name),
				When:    time.Now(),
			}
			h.broadcast <- m
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				cm := websocket.FormatCloseMessage(websocket.CloseNormalClosure, fmt.Sprintf("closing connection for client: ", client.name))
				if err := client.conn.WriteMessage(websocket.CloseMessage, cm); err != nil {
					// handle error
				}
				err := client.conn.Close()
				if err != nil {
					fmt.Println("Error in closing client", client, err)
					// TODO: re think cleanup here!
					return
				}
				delete(h.clients, client)
				close(client.send)
				m := Message{
					Name:    "",
					Message: fmt.Sprintf("%s has left chat", client.name),
					When:    time.Now(),
				}
				h.broadcast <- m
			}

		case message := <-h.broadcast:
			fmt.Printf("Broadcasting message. %v\n", message)
			// send Message to redis
			h.redisConnector.client.Publish(ctx, CHATROOM, message)

			// send to local clients
			for client := range h.clients {
				if client.name != message.Name {
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
