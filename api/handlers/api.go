package handlers

import (
	"api/service"
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
	Storage   service.Storage
	DB        *sql.DB
	JwtSecret string

	WsClients WsClient
	Upgrader  websocket.Upgrader
}

func (A *API) secret() []byte {
	return []byte(A.JwtSecret)
}

func New(db *sql.DB) *API {
	return &API{
		DB:        db,
		JwtSecret: "secret-secret",

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
