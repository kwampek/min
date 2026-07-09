package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type WSHandlerFunc func(userID int, payload json.RawMessage) error

type WsRouter struct {
	handlers map[string]WSHandlerFunc
}

func NewWsRouter(api *API) WsRouter {

	return WsRouter{
		handlers: map[string]WSHandlerFunc{
			"search":                api.SearchUsersHandler,
			"new_message":           api.AddMessageHandler,
			"create_chat":           api.CreatePrivateChatHandler,
			"create_folder":         api.CreateFolderHandler,
			"edit_profile":          api.EditProfileHandler,
			"edit_chat_name":        api.EditChatNameHandler,
			"edit_folder_name":      api.EditFolderNameHandler,
			"toggle_chat_in_folder": api.ToggleChatInFolderHandler,
		},
	}
}

func (r *WsRouter) Get(name string) (WSHandlerFunc, bool) {
	handler, ok := r.handlers[name]
	return handler, ok
}

type WSMessageWrapper struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

func (A *API) authenticateWS(r *http.Request) (int, error) {
	tokenString := r.URL.Query().Get("token")
	if tokenString == "" {
		return 0, fmt.Errorf("missing token")
	}

	claims, err := A.AuthService.Jwt.ValidateToken(tokenString)
	if err != nil {
		return 0, err
	}

	exists, err := A.AuthService.Storage.IsTokenExists(tokenString)
	if err != nil {
		return 0, err
	}

	if !exists {
		return 0, fmt.Errorf("token revoked")
	}

	// may be better to send full claims
	return claims.UserID, nil
}

func (A *API) sendInitialState(userID int, conn *websocket.Conn) error {
	allMessages, err := A.MessageService.LoadAllMessages(userID)
	if err != nil {
		log.Println("load messages error:", err)
		_ = A.WsService.FastSend(conn, map[string]interface{}{
			"type":  "error",
			"error": "failed to load messages",
		})
		return err
	}

	return A.WsService.FastSend(conn, map[string]interface{}{
		"type":          "initial_state",
		"initial_state": allMessages,
	})
}

func (A *API) WsHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := A.authenticateWS(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	conn, err := A.WsService.Upgrade(w, r)
	if err != nil {
		return
	}

	A.WsService.AddClient(userID, conn)
	defer A.WsService.RemoveClient(userID)

	if err := A.sendInitialState(userID, conn); err != nil {
		return
	}

	for {

		_, data, err := conn.ReadMessage()
		if err != nil {
			return
		}

		var msg WSMessageWrapper

		if err := json.Unmarshal(data, &msg); err != nil {
			continue
		}

		handler, ok := A.WsRouter.Get(msg.Type)
		if !ok {
			continue
		}

		if err := handler(userID, msg.Payload); err != nil {
			log.Println(err)
		}
	}
}
