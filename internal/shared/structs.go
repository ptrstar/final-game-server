package shared

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type ClientRequest struct {
	W     http.ResponseWriter
	R     *http.Request
	GType string
	RID   string
	Conn  *websocket.Conn
}

type ShareableGameState struct {
	Type        string    `json:"type"`
	Id          string    `json:"id"`
	Capacity    int       `json:"capacity"`
	Status      string    `json:"status"`
	CanJoin     bool      `json:"can_join"`
	PlayerCount int       `json:"player_count"`
	CreatedAt   time.Time `json:"created_at"`
	Custom      []byte    `json:"custom"` // See note below on []byte
}

type Vec2 struct {
	X float32
	Y float32
}
type Col struct {
	R uint8
	G uint8
	B uint8
}
