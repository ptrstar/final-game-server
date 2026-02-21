package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
		//origin := r.Header.Get("Origin")
		//return origin == "https://its.trojanos.ch" || origin == "https://play.trojanos.ch" || origin == ""
	},
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow its.. to fetch data
		w.Header().Set("Access-Control-Allow-Origin", "https://its.trojanos.ch")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle browser pre-flight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()

	hub := NewHub()

	mux.Handle("/", http.FileServer(http.Dir("./static")))

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
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

	wrappedMux := enableCORS(mux)

	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", wrappedMux)
}
