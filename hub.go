package main

import (
	"errors"
	"final-game-server/internal/engine"
	"final-game-server/internal/games/info"
	"final-game-server/internal/games/paint"
	"final-game-server/internal/shared"
	"fmt"
	"math/rand/v2"
	"time"
)

var validGameTypes = map[string]struct{}{
	"info":  {},
	"paint": {},
}

type Hub struct {
	Rooms map[string]map[string]*engine.Room
	Join  chan *shared.ClientRequest
	// mu     sync.Mutex
	nextId int
}

func NewHub() *Hub {
	hub := &Hub{
		Rooms: make(map[string]map[string]*engine.Room),
		Join:  make(chan *shared.ClientRequest),
		// mu:     sync.Mutex{},
		nextId: 0,
	}

	hub.Rooms["info"] = make(map[string]*engine.Room)
	room := engine.NewRoom("all", "info", 128)
	room.Game = info.NewInfo(room)
	hub.Rooms["info"]["all"] = room

	return hub
}

func (h *Hub) ConnectClientToRoom(client *engine.Client, cr *shared.ClientRequest) error {
	// h.mu.Lock()
	// defer h.mu.Unlock()

	if cr.GType == "info" {
		room := h.Rooms["info"]["all"]
		return room.AddClient(client)
	}

	if _, ok := validGameTypes[cr.GType]; !ok {
		return errors.New("invalid game type")
	}

	if h.Rooms[cr.GType] == nil {
		h.Rooms[cr.GType] = make(map[string]*engine.Room)
	}

	var targetRoom *engine.Room

	if cr.RID != "" && cr.RID != "any" {
		if r, exists := h.Rooms[cr.GType][cr.RID]; exists {
			if err := r.AddClient(client); err == nil {
				targetRoom = r
			}
		}
	}

	if targetRoom == nil {
		for _, r := range h.Rooms[cr.GType] {
			if err := r.AddClient(client); err == nil {
				targetRoom = r
				break
			}
		}
	}

	if targetRoom == nil {
		newId := fmt.Sprintf("%06x", rand.Uint32())[:6]
		targetRoom = h.CreateRoom(cr.GType, newId)

		h.Rooms[cr.GType][newId] = targetRoom
		if err := targetRoom.AddClient(client); err != nil {
			return err
		}
	}

	return nil
}

func (h *Hub) CreateRoom(GType string, Id string) *engine.Room {
	if GType == "paint" {
		room := engine.NewRoom(Id, GType, 4)
		room.Game = paint.NewPaintGame(room)
		return room
	}
	panic("Unable to find Roomconstructor")
}

func (h *Hub) Run() {
	hubTicker := time.NewTicker(1000 * time.Millisecond)

	for {
		select {
		case cr := <-h.Join:

			client := engine.NewClient(h.getNextId(), cr.Conn)

			h.ConnectClientToRoom(client, cr)

		case <-hubTicker.C:
			GSList := make([]*shared.ShareableGameState, 0)
			for Type := range h.Rooms {
				for _, room := range h.Rooms[Type] {
					GSList = append(GSList, room.Game.GetShareableGameState())
				}
			}
			if infoGame, ok := h.Rooms["info"]["all"].Game.(*info.Info); ok {
				infoGame.SetGSList(GSList)
			}
		}

	}
}

func (h *Hub) getNextId() int {
	h.nextId++
	return h.nextId
}
