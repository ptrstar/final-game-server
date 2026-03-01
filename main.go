package main

import (
	"final-game-server/internal/shared"
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
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")

		// Only force the JS content-type if the URL actually ends in .js
		if len(r.URL.Path) >= 3 && r.URL.Path[len(r.URL.Path)-3:] == ".js" {
			w.Header().Set("Content-Type", "application/javascript")
		}

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
		rID := r.URL.Query().Get("id")

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		cr := &shared.ClientRequest{
			W:     w,
			R:     r,
			GType: gType,
			RID:   rID,
			Conn:  conn,
		}

		hub.Join <- cr
	})

	wrappedMux := enableCORS(mux)

	go hub.Run()

	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", wrappedMux)

}
