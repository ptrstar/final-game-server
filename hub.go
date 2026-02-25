package main

import (
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

var validGameTypes = map[string]struct{}{
	"wasd": {},
}

type Hub struct {
	Games map[string]map[string]*Room
	Info  map[string]*Room
	Join  chan *ClientRequest
	mu    sync.Mutex
}

func NewHub() *Hub {
	hub := &Hub{
		Games: make(map[string]map[string]*Room),
		Info:  make(map[string]*Room),
		Join:  make(chan *ClientRequest),
		mu:    sync.Mutex{},
	}

	hub.Info["status"] = NewRoom("status", "hub", 1000)
	hub.Info["info"] = NewRoom("info", "hub", 1000)

	return hub
}

func (h *Hub) ConnectClientToRoom(gameType, roomID string, clientRequest ClientRequest) *Room {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.Games[gameType] == nil {
		h.Games[gameType] = make(map[string]*Room)
	}

	room, exists := h.Games[gameType][roomID]
	if (!exists) || (len(room.Clients) >= room.Game.MaxPlayers) {
		foundRoom := false
		for _, existingRoom := range h.Games[gameType] {
			if existingRoom.AddPlayer(client) {
				room = existingRoom
				roomID = existingRoom.ID
				foundRoom = true
				break
			}
		}
		if !foundRoom {
			if exists {
				roomID = uuid.NewString()
			}
			room = NewRoom(roomID, gameType)
			h.Games[gameType][roomID] = room
			go room.Run()
		}
	}
	// return room
}

func (h *Hub) Run() {
	broadcastTicker := time.NewTicker(50 * time.Millisecond)

	for {
		select {
		case clientRequest := <-h.Join:
			// connected to hub

			conn, err := upgrader.Upgrade(clientRequest.w, &clientRequest.r, nil)
			if err != nil {
				return
			}

			if clientRequest.hInfo == "info" || clientRequest.hInfo == "status" {
				room := h.Info[clientRequest.hInfo]
				client := &Client{Conn: conn, Send: make(chan []byte, 256), Room: room}
				if room.AddPlayer(client) != nil {
					conn.Close()
					log.Printf("Failed to connect client to hub")
				} else {
					go client.WritePump()
					client.ReadPump()
				}
				return
			}

			_, ok := validGameTypes[clientRequest.gType]
			if !ok {
				log.Printf("Clientconnection failed due to invalid game type")
				conn.Close()
				return
			}

			client := &Client{Conn: conn, Send: make(chan []byte, 256), Room: nil}
			// TODO: figure out how to send clientrequest to this function which should use game.AddPlayer to dynamically init the correct Player from clientrequest data
			h.ConnectClientToRoom(clientRequest.gType, clientRequest.rID, client)
			go client.WritePump()
			client.ReadPump()

			// room := hub.GetOrCreateRoom(gType, rID)
			// conn, err := upgrader.Upgrade(w, r, nil)
			// if err != nil {
			// 	return
			// }

			// use clientconsturction NewClient instead!!!!!!!
			// client := &Client{Conn: conn, Send: make(chan []byte, 256), Room: room}
			// room.Join <- client

			// go client.WritePump()
			// client.ReadPump()

		case <-broadcastTicker.C:

		}

	}
}
