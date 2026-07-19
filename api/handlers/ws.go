package handlers

import (
	"api/models"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type WSHandlerFunc func(id models.Identifier, payload json.RawMessage) error

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
			"logout":                api.LogoutHandler,
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

func (a *API) authenticateWS(r *http.Request) (models.Identifier, error) {
	token := r.URL.Query().Get("token")
	if token == "" {
		return models.Identifier{}, errors.New("missing token")
	}

	fmt.Println("AUTH TRY: ", token)
	return a.AuthService.ValidateSession(token)
}

func (A *API) sendInitialState(ID models.Identifier, conn *websocket.Conn) error {
	allMessages, err := A.MessageService.LoadAllMessages(ID.UserID)
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
	ID, err := A.authenticateWS(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	conn, err := A.WsService.Upgrade(w, r)
	if err != nil {
		return
	}

	A.WsService.AddClient(ID, conn)
	defer A.WsService.RemoveConnection(ID)

	if err := A.sendInitialState(ID, conn); err != nil {
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

		if err := handler(ID, msg.Payload); err != nil {
			log.Println(err)
		}
	}
}
