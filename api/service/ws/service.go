package ws

import (
	"api/models"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type WsService struct {
	mu           sync.RWMutex
	clients      map[int]*websocket.Conn
	userSessions map[int][]int

	upgrader websocket.Upgrader
}

func NewWsService() WsService {
	return WsService{
		mu:           sync.RWMutex{},
		clients:      make(map[int]*websocket.Conn),
		userSessions: make(map[int][]int),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (ws *WsService) AddClient(ID models.Identifier, conn *websocket.Conn) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	ws.clients[ID.SessionID] = conn
	ws.userSessions[ID.UserID] = append(ws.userSessions[ID.UserID], ID.SessionID)
}

func (ws *WsService) RemoveConnection(ID models.Identifier) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	if conn, ok := ws.clients[ID.SessionID]; ok {
		conn.Close()
		delete(ws.clients, ID.SessionID)
	}
}

func (ws *WsService) RemoveClient(ID models.Identifier) {
	ws.RemoveConnection(ID)

	// std::swap with last removing
	userSessions := ws.userSessions[ID.UserID]
	for i, sid := range userSessions {
		if sid == ID.SessionID {
			userSessions[i] = userSessions[len(userSessions)-1]
			break
		}
	}
	ws.userSessions[ID.UserID] = userSessions[:len(userSessions)-1]
}

func (h *WsService) FastSend(conn *websocket.Conn, payload any) error {
	if conn == nil {
		return nil
	}

	return conn.WriteJSON(payload)
}

func (h *WsService) Send(ID models.Identifier, payload any) error {
	h.mu.RLock()
	conn, ok := h.clients[ID.SessionID]
	h.mu.RUnlock()

	if !ok || conn == nil {
		return nil
	}

	err := conn.WriteJSON(payload)
	if err != nil {
		h.RemoveConnection(ID)
	}

	return err
}

func (h *WsService) BroadcastToUsers(users []int, senderID int, payload any) error {
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

func (ws *WsService) Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	return ws.upgrader.Upgrade(w, r, nil)
}

func (ws *WsService) Kick(userID int) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	if conn, ok := ws.clients[userID]; ok {
		conn.WriteControl(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(
				websocket.ClosePolicyViolation,
				"logged out",
			),
			time.Now().Add(time.Second),
		)

		conn.Close()

		delete(ws.clients, userID)
	}
}
