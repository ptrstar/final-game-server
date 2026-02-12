package main

import (
	"final-game-server/game/wasd"
	"log"
	"time"
)

type ClientInput struct {
	Client *Client
	Data   []byte
}

type Room struct {
	ID      string
	Type    string
	Clients map[*Client]bool
	Join    chan *Client
	Leave   chan *Client
	Input   chan ClientInput
	Game    *wasd.WASDGame // Expandable to an interface for multiple types
}

func NewRoom(id, gType string) *Room {
	return &Room{
		ID:      id,
		Type:    gType,
		Clients: make(map[*Client]bool),
		Join:    make(chan *Client),
		Leave:   make(chan *Client),
		Input:   make(chan ClientInput),
		Game:    wasd.NewGame(3),
	}
}

func (r *Room) Broadcast(state []byte) {
	for c := range r.Clients {
		select {
		case c.Send <- state:
		default:
			// Handle slow clients
		}
	}
}

func (r *Room) Run() {
	physicsTicker := time.NewTicker(16 * time.Millisecond)
	broadcastTicker := time.NewTicker(50 * time.Millisecond)
	defer physicsTicker.Stop()
	defer broadcastTicker.Stop()

	for {
		select {
		case c := <-r.Join:
			r.Clients[c] = true
			r.Game.AddPlayer(c)
			log.Printf("Client %d joined room %s\n", r.Game.Players[c].ID, r.ID)
		case c := <-r.Leave:
			log.Printf("Client %d left room %s\n", r.Game.Players[c].ID, r.ID)
			delete(r.Clients, c)
			r.Game.RemovePlayer(c)
		case input := <-r.Input:
			r.Game.HandleInput(input.Client, input.Data)
		case <-broadcastTicker.C:
			if len(r.Clients) > 0 {
				state := r.Game.SerializeState()
				r.Broadcast(state)
			}
		case <-physicsTicker.C:
			r.Game.Update()
		}
	}
}
