package main

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn *websocket.Conn
	Send chan []byte
	Room *Room
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
		c.Room.Input <- ClientInput{Client: c, Data: message}
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
