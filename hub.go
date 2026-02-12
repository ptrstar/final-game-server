package main

import (
	"sync"

	"github.com/google/uuid"
)

type Hub struct {
	Games map[string]map[string]*Room
	mu    sync.Mutex
}

func NewHub() *Hub {
	return &Hub{Games: make(map[string]map[string]*Room)}
}

func (h *Hub) GetOrCreateRoom(gameType, roomID string) *Room {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.Games[gameType] == nil {
		h.Games[gameType] = make(map[string]*Room)
	}

	room, exists := h.Games[gameType][roomID]
	if (!exists) || (len(room.Clients) >= room.Game.MaxPlayers) {
		foundRoom := false
		for _, existingRoom := range h.Games[gameType] {
			if len(existingRoom.Clients) < existingRoom.Game.MaxPlayers {
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
	return room
}
