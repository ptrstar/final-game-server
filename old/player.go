package main

import (
	"log"

	"github.com/gorilla/websocket"
)

type PlayerInterface interface {
	ReadPump()
	WritePump()
	PutMessage(msg []byte)
}

type Client struct {
	ID string

	Conn *websocket.Conn
	Send chan []byte
	Room *Room
}

func NewClient(conn *websocket.Conn, room *Room) *Client {
	return &Client{Conn: conn, Send: make(chan []byte, 256), Room: room}
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
		// Pass irregular client updates to the room's input channel
		c.Room.Input <- ClientInput{Player: c, Data: message}
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
