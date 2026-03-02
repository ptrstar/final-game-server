package engine

import (
	"errors"
	"log"
	"sync"
	"time"
)

type Room struct {
	mu          sync.RWMutex
	Id          string
	Type        string
	Clients     map[*Client]*Client
	Leave       chan *Client
	Input       chan *ClientInput
	Game        GameItf
	CreatedAt   time.Time
	UpdateCycle int

	Capacity int
	nextId   int
}

func NewRoom(Id string, Type string, Capacity int) *Room {

	room := &Room{
		mu:          sync.RWMutex{},
		Id:          Id,
		Type:        Type,
		Clients:     make(map[*Client]*Client),
		Leave:       make(chan *Client),
		Input:       make(chan *ClientInput),
		CreatedAt:   time.Now(),
		UpdateCycle: 1000,

		Capacity: Capacity,
		nextId:   0,
	}

	go room.Run()

	return room
}

func (r *Room) Run() {
	broadcastTicker := time.NewTicker(time.Duration(r.UpdateCycle) * time.Millisecond)

	for {
		select {
		case <-broadcastTicker.C:
			state := r.Game.SerializeState()
			r.Broadcast(state)
		case client := <-r.Leave:
			r.RemoveClient(client)
		case input := <-r.Input:
			r.Game.HandleInput(input)
		}
	}
}
func (r *Room) Broadcast(data []byte) {
	for p := range r.Clients {
		p.PutMessage(data)
	}
}
func (r *Room) AddClient(client *Client) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.Clients) >= r.Capacity {
		return errors.New("Room is full")
	}

	r.Clients[client] = client
	client.SetRoom(r)
	r.Game.AddPlayer(client)

	go client.WritePump()
	go client.ReadPump()

	log.Printf("Client %d connected to Room[\"%s\"][\"%s\"]\n", client.Id, r.Type, r.Id)

	return nil
}
func (r *Room) RemoveClient(client *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.Game.RemovePlayer(client)

	log.Printf("Client %d disconnected from Room[\"%s\"][\"%s\"]\n", client.Id, r.Type, r.Id)

	delete(r.Clients, client)
}
func (r *Room) HandleInput(client *Client, input []byte) {}
