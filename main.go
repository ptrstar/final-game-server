package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type ClientRequest struct {
	w     http.ResponseWriter
	r     http.Request
	gType string
	rID   string
	hInfo string
	conn  *websocket.Conn
}

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
		h := r.URL.Query().Get("hub")

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		clientRequest := &ClientRequest{
			w:     w,
			r:     *r,
			gType: gType,
			rID:   rID,
			hInfo: h,
			conn:  conn,
		}
		hub.Join <- clientRequest
	})

	wrappedMux := enableCORS(mux)

	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", wrappedMux)

	hub.Run()
}
