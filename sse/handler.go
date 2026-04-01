// sse/handler.go
package sse

import (
	"fmt"
	"net/http"
	"task-tracker-api/utils"

	"github.com/google/uuid"
)

type SSEHandler struct {
	hub *Hub
}

func NewSSEHandler(hub *Hub) *SSEHandler {
	return &SSEHandler{hub: hub}
}

func (h *SSEHandler) Stream(w http.ResponseWriter, r *http.Request) error {
	// get boardID from query instead of teamID
	boardID := r.URL.Query().Get("board_id")
	if boardID == "" {
		return utils.BadRequest("board_id is required")
	}

	claims := r.Context().Value(utils.ClaimsKey).(*utils.Claims)

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	client := &Client{
		ID:      fmt.Sprintf("%s-%s", claims.UserID, uuid.New().String()),
		Channel: make(chan string, 10),
		BoardID: boardID, // use boardID here
	}

	h.hub.AddClient(client)
	defer h.hub.RemoveClient(client.ID)

	fmt.Fprintf(w, "event: connected\ndata: {\"message\": \"connected to board %s\"}\n\n", boardID)
	w.(http.Flusher).Flush()

	for {
		select {
		case message, ok := <-client.Channel:
			if !ok {
				return nil
			}
			fmt.Fprint(w, message)
			w.(http.Flusher).Flush()

		case <-r.Context().Done():
			return nil
		}
	}
}
