package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func main() {
	hub := NewHub()

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.Handle("/snippets/", http.StripPrefix("/snippets/", http.FileServer(http.Dir("./snippets"))))

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		gType := r.URL.Query().Get("type")
		rID := r.URL.Query().Get("room")
		if gType == "" || rID == "" {
			return
		}

		room := hub.GetOrCreateRoom(gType, rID)
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		client := &Client{Conn: conn, Send: make(chan []byte, 256), Room: room}
		room.Join <- client

		go client.WritePump()
		client.ReadPump()
	})

	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}
