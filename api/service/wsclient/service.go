package wsclient

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type WsClient struct {
	mu      sync.RWMutex
	clients map[int]*websocket.Conn
}

func NewWsClient() WsClient {
	return WsClient{
		mu:      sync.RWMutex{},
		clients: make(map[int]*websocket.Conn),
	}
}

func (h *WsClient) Send(userID int, payload any) error {
	h.mu.RLock()
	conn, ok := h.clients[userID]
	h.mu.RUnlock()

	if !ok || conn == nil {
		return nil
	}

	return conn.WriteJSON(payload)
}

func (h *WsClient) BroadcastToUsers(users []int, senderID int, payload any) error {
	for _, uid := range users {
		if uid == senderID {
			continue
		}
		conn, ok := h.clients[uid]
		if !ok || conn == nil {
			continue
		}

		if err := conn.WriteJSON(payload); err != nil {
			log.Printf("WS send error to user %d: %v", uid, err)
			continue
		}
	}

	// TODO - maybe need to add errors

	return nil
}
