// Inspiration and help taken from the gorilla websocket examples

package natsterui

import (
	log "log/slog"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type SocksHub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewHub() *SocksHub {
	return &SocksHub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *SocksHub) BroadcastMessage(msg []byte) {
	h.broadcast <- msg
}

func (h *SocksHub) Run() {
	log.Info("Starting Websocket Hub")
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				client.send <- message
			}
		}
	}
}
