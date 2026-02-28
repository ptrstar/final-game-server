package main

import (
	"final-game-server/internal/engine"
	"final-game-server/internal/games/info"
	"final-game-server/internal/shared"
	"log"
	"sync"
	"time"
)

type Hub struct {
	Rooms map[string]map[string]*engine.Room
	Join  chan *shared.ClientRequest
	mu    sync.Mutex
}

func NewHub() *Hub {
	hub := &Hub{
		Rooms: make(map[string]map[string]*engine.Room),
		Join:  make(chan *shared.ClientRequest),
		mu:    sync.Mutex{},
	}

	hub.Rooms["info"] = make(map[string]*engine.Room)
	room, _ := engine.NewRoom("All", "info", 1000)
	room.Game = info.NewInfo()
	hub.Rooms["info"]["All"] = room

	return hub
}

func (h *Hub) ConnectClientToRoom(client *engine.Client, cr *shared.ClientRequest) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if cr.GType == "Info" {
		room := h.Rooms["Info"]["ALL"]
		if room.AddClient(client) != nil {
			log.Println("Unable to add client to info room")
		}
	}

	// if h.Games[gameType] == nil {
	// 	h.Games[gameType] = make(map[string]*Room)
	// }

	// room, exists := h.Games[gameType][roomID]
	// if (!exists) || (len(room.Clients) >= room.Game.MaxPlayers) {
	// 	foundRoom := false
	// 	for _, existingRoom := range h.Games[gameType] {
	// 		if existingRoom.AddPlayer(client) {
	// 			room = existingRoom
	// 			roomID = existingRoom.ID
	// 			foundRoom = true
	// 			break
	// 		}
	// 	}
	// 	if !foundRoom {
	// 		if exists {
	// 			roomID = uuid.NewString()
	// 		}
	// 		room = NewRoom(roomID, gameType)
	// 		h.Games[gameType][roomID] = room
	// 		go room.Run()
	// 	}
	// }
	// // return room
}

func (h *Hub) Run() {
	broadcastTicker := time.NewTicker(50 * time.Millisecond)

	for {
		select {
		case cr := <-h.Join:

			conn, err := upgrader.Upgrade(cr.W, &cr.R, nil)
			if err != nil {
				return
			}

			client := engine.NewClient(conn)

			h.ConnectClientToRoom(client, cr)

			// 		if clientRequest.HInfo == "info" || clientRequest.HInfo == "status" {
			// 			room := h.Info[clientRequest.HInfo]
			// 			client := &Client{Conn: conn, Send: make(chan []byte, 256), Room: room}
			// 			if room.AddPlayer(client) != nil {
			// 				conn.Close()
			// 				log.Printf("Failed to connect client to hub")
			// 			} else {
			// 				go client.WritePump()
			// 				client.ReadPump()
			// 			}
			// 			return
			// 		}

			// 		_, ok := validGameTypes[clientRequest.GType]
			// 		if !ok {
			// 			log.Printf("Clientconnection failed due to invalid game type")
			// 			conn.Close()
			// 			return
			// 		}

			// 		client := &Client{Conn: conn, Send: make(chan []byte, 256), Room: nil}
			// 		// TODO: figure out how to send clientrequest to this function which should use game.AddPlayer to dynamically init the correct Player from clientrequest data
			// 		h.ConnectClientToRoom(clientRequest.GType, clientRequest.RID, client)
			// 		go client.WritePump()
			// 		client.ReadPump()

			// 		// room := hub.GetOrCreateRoom(gType, rID)
			// 		// conn, err := upgrader.Upgrade(w, r, nil)
			// 		// if err != nil {
			// 		// 	return
			// 		// }

			// 		// use clientconsturction NewClient instead!!!!!!!
			// 		// client := &Client{Conn: conn, Send: make(chan []byte, 256), Room: room}
			// 		// room.Join <- client

			// 		// go client.WritePump()
			// 		// client.ReadPump()

		case <-broadcastTicker.C:

		}

	}
}
