package engine

import (
	"errors"
	"sync"
)

var validGameTypes = map[string]struct{}{
	"info": {},
}

type Room struct {
	mu      sync.RWMutex
	Id      string
	Type    string
	Clients map[*Client]*Client
	Leave   chan *Client
	Input   chan *ClientInput
	Game    GameItf

	Capacity int
	nextId   int
}

func NewRoom(Id string, Type string, Capacity int) (*Room, error) {

	_, ok := validGameTypes[Type]
	if !ok {
		return nil, errors.New("Invalid RoomType")
	}

	room := &Room{
		mu:      sync.RWMutex{},
		Id:      Id,
		Type:    Type,
		Clients: make(map[*Client]*Client),
		Leave:   make(chan *Client),
		Input:   make(chan *ClientInput, 0),

		Capacity: Capacity,
		nextId:   0,
	}

	go room.Run()

	return room, nil
}

func (r *Room) Run()                  {}
func (r *Room) Broadcast(data []byte) {}
func (r *Room) AddClient(client *Client) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.Clients[client] = client
	client.SetRoom(r)

	go client.WritePump()
	go client.ReadPump()

	return nil
}
func (r *Room) RemoveClient(client *Client)              {}
func (r *Room) HandleInput(client *Client, input []byte) {}
