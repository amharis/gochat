// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"time"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan message

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan message, 1024),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			h.broadcast <- message{
				Name:    "",
				Message: fmt.Sprintf("%s has joined chat", client.name),
				When:    time.Now(),
			}
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				cm := websocket.FormatCloseMessage(websocket.CloseNormalClosure, fmt.Sprintf("closing connection for client: ", client.name))
				if err := client.conn.WriteMessage(websocket.CloseMessage, cm); err != nil {
					// handle error
				}
				client.conn.Close()
				delete(h.clients, client)
				close(client.send)
			}
			h.broadcast <- message{
				Name:    "",
				Message: fmt.Sprintf("%s has left chat", client.name),
				When:    time.Now(),
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
