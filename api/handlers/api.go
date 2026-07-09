package handlers

import (
	"api/service/auth"
	"api/service/chats"
	"api/service/jwt"
	"api/service/media"
	"api/service/messages"
	"api/service/wsclient"
	"api/storage"
	"database/sql"
	"net/http"

	"github.com/gorilla/websocket"
)

type API struct {
	DB *sql.DB

	AuthService    auth.Service
	MessageService messages.MessageService
	MediaService   media.MediaService
	ChatService    chats.ChatService

	WsClients wsclient.WsClient
	Upgrader  websocket.Upgrader
}

func New(db *sql.DB) *API {
	secret := "secret-secret"
	issuer := "auth.min.com"
	expiryAt := 7

	JwtService := jwt.NewService(secret, issuer, expiryAt)
	Storage := storage.NewStorage(db)

	return &API{
		DB:             db, // maybe obsolete
		AuthService:    auth.NewService(Storage, JwtService),
		MessageService: messages.NewMessagesService(Storage),
		MediaService:   media.NewMediaService(Storage),
		ChatService:    chats.NewChatService(Storage),

		WsClients: wsclient.NewWsClient(),
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}
