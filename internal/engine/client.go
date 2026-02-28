package engine

import (
	"log"

	"github.com/gorilla/websocket"
)

type ClientInput struct {
	Client *Client
	Data   []byte
}

type Client struct {
	Conn *websocket.Conn
	Send chan []byte
	Room *Room
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{Conn: conn, Send: make(chan []byte, 256)}
}

func (c *Client) SetRoom(room *Room) {
	c.Room = room
}

func (c *Client) ReadPump() {
	defer func() {
		c.Room.Leave <- c
		c.Conn.Close()
	}()
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		c.Room.Input <- &ClientInput{Client: c, Data: message}
	}
}

func (c *Client) WritePump() {
	for msg := range c.Send {
		if err := c.Conn.WriteMessage(websocket.BinaryMessage, msg); err != nil {
			log.Println("Write error:", err)
			return
		}
	}
}

// called by roomthread
func (c *Client) PutMessage(msg []byte) {
	select {
	case c.Send <- msg:
	default:
		// optional: drop message if buffer is full to prevent hanging the room
	}
}
