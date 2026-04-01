package sse

import (
	"fmt"
	"sync"
)

type Client struct {
	ID      string
	Channel chan string
	BoardID string // renamed from TeamID
}

type Hub struct {
	clients map[string]*Client
	mu      sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]*Client),
	}
}

func (h *Hub) AddClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[client.ID] = client
}

func (h *Hub) RemoveClient(clientID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if client, exists := h.clients[clientID]; exists {
		close(client.Channel)
		delete(h.clients, clientID)
	}
}

// renamed from BroadcastToTeam — now broadcasts by boardID
func (h *Hub) Broadcast(boardID string, event string, data string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	message := fmt.Sprintf("event: %s\ndata: %s\n\n", event, data)

	for _, client := range h.clients {
		if client.BoardID == boardID { // match by boardID
			select {
			case client.Channel <- message:
			default:
			}
		}
	}
}
