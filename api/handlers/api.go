package handlers

import (
	"api/service/auth"
	"api/service/chats"
	"api/service/folders"
	"api/service/media"
	"api/service/messages"
	"api/service/tokens"
	"api/service/users"
	"api/service/ws"
	"api/storage"
	"database/sql"
)

type API struct {
	DB *sql.DB

	AuthService    auth.AuthService
	MessageService messages.MessageService
	MediaService   media.MediaService
	ChatService    chats.ChatService
	FolderService  folders.FolderService
	UserService    users.UserService

	WsService ws.WsService
	WsRouter  WsRouter
}

func New(db *sql.DB) *API {

	storage := storage.NewStorage(db)

	api := &API{
		DB: db,

		AuthService: auth.NewAuthService(
			storage,
			tokens.NewTokenService(
				"secret-secret",
				"auth.min.com",
				7,
			),
		),
		MessageService: messages.NewMessagesService(storage),
		MediaService:   media.NewMediaService(storage),
		ChatService:    chats.NewChatService(storage),
		FolderService:  folders.NewFolderService(storage),
		UserService:    users.NewUserService(storage),

		WsService: ws.NewWsService(),
	}

	api.WsRouter = NewWsRouter(api)

	return api
}
