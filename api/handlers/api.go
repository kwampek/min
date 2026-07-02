package handlers

import (
	"api/service/auth"
	"api/service/jwt"
	"api/storage"
	"database/sql"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type WsClient struct {
	sync.RWMutex // Write with lock, Read free
	Conns        map[int]*websocket.Conn
}

type API struct {
	DB *sql.DB

	AuthService auth.Service

	WsClients WsClient
	Upgrader  websocket.Upgrader
}

func New(db *sql.DB) *API {
	secret := "secret-secret"
	issuer := "auth.min.com"
	expiryAt := 7

	JwtService := jwt.NewService(secret, issuer, expiryAt)
	Storage := storage.NewStorage(db)

	return &API{
		DB:          db, // maybe obsolete
		AuthService: auth.NewService(Storage, JwtService),

		WsClients: WsClient{
			Conns: make(map[int]*websocket.Conn),
		},
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}
