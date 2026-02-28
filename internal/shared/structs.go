package shared

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type ClientRequest struct {
	W     http.ResponseWriter
	R     http.Request
	GType string
	RID   string
	Conn  *websocket.Conn
}
