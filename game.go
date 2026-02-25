package main

import (
	"errors"
	"log"
	"math/rand/v2"
	"sync"
	"time"
)

// TODO: figure this out
type ClientInput struct {
	Player PlayerInterface
	Data   []byte
}

// TODO: figure this out too
type PlayerState struct {
	ID       int
	Position Vector2
	Color    RGBColor
	Keyboard map[byte]bool
}

type GameInterface interface {
	Run()
	Broadcast([]byte)
	// Obligatory game functions
	AddPlayer(PlayerInterface)
	RemovePlayer(PlayerInterface)
	HandleInput(PlayerInterface)
	GameUpdate()
	SerializeState()
}

// Embedding of the GameInterface
type Room struct {
	mu      sync.RWMutex
	ID      string
	Type    string
	Players map[PlayerInterface]*PlayerState
	Leave   chan PlayerInterface
	Input   chan ClientInput

	MaxPlayers int
	nextID     int
}

func NewRoom(id, gType string, maxClients int) *Room {
	return &Room{
		ID:      id,
		Type:    gType,
		Players: make(map[PlayerInterface]*PlayerState),
		Leave:   make(chan PlayerInterface),
		Input:   make(chan ClientInput),
	}
}

func (r *Room) Broadcast(state []byte) {
	for p := range r.Players {
		p.PutMessage(state)
	}
}

func (r *Room) Run() {
	physicsTicker := time.NewTicker(16 * time.Millisecond)
	broadcastTicker := time.NewTicker(50 * time.Millisecond)
	defer physicsTicker.Stop()
	defer broadcastTicker.Stop()

	for {
		select {
		case p := <-r.Leave:
			r.RemovePlayer(p)
		case input := <-r.Input:
			r.HandleInput(input.Player, input.Data)
		case <-broadcastTicker.C:
			// TODO: check if lock is needed
			if len(r.Players) > 0 {
				state := r.SerializeState()
				r.Broadcast(state)
			}
		case <-physicsTicker.C:
			r.GameUpdate()
		}
	}
}

// called by hub thread
func (g *Room) AddPlayer(c PlayerInterface) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if len(g.Players) >= g.MaxPlayers {
		return errors.New("rejected")
	}

	g.Players[c] = &PlayerState{ID: g.nextID, Color: RGBColor{R: uint8(rand.IntN(255)), G: uint8(rand.IntN(255)), B: uint8(rand.IntN(255))}, Position: Vector2{X: 100, Y: 100}, Keyboard: make(map[byte]bool)}
	g.nextID++

	log.Printf("Client %d joined room %s\n", g.Players[c].ID, g.ID)
	return nil
}
func (g *Room) RemovePlayer(c PlayerInterface) {
	g.mu.Lock()
	defer g.mu.Unlock()

	// TODO: use channel to tell hub to remove game if playercount drops to 0

	log.Printf("Client %d left room %s\n", g.Players[c].ID, g.ID)
	delete(g.Players, c)

}
func (g *Room) HandleInput(c PlayerInterface, data []byte) {

	// if len(data) == 0 {
	// 	return
	// }
	// g.mu.Lock()
	// defer g.mu.Unlock()
	// p, ok := g.Players[c]
	// if !ok {
	// 	return
	// }

	// key := data[0] >> 1
	// var value bool = (data[0] & 0b01) != 0

	// p.keyboard[key] = value
}
func (g *Room) GameUpdate() {
	// for _, player := range g.Players {
	// 	if player.keyboard[0] {
	// 		player.Y -= 5
	// 	}
	// 	if player.keyboard[1] {
	// 		player.X -= 5
	// 	}
	// 	if player.keyboard[2] {
	// 		player.Y += 5
	// 	}
	// 	if player.keyboard[3] {
	// 		player.X += 5
	// 	}
	// }
}
func (g *Room) SerializeState() []byte {

	buf := make([]byte, 1+(len(g.Players)*11))
	// buf[0] = byte(len(g.Players))

	// i := 1
	// for _, p := range g.Players {
	// 	buf[i] = p.ID
	// 	var keyboard uint16 = 0
	// 	if p.keyboard[0] {
	// 		keyboard |= 1
	// 	}
	// 	if p.keyboard[1] {
	// 		keyboard |= 2
	// 	}
	// 	if p.keyboard[2] {
	// 		keyboard |= 4
	// 	}
	// 	if p.keyboard[3] {
	// 		keyboard |= 8
	// 	}
	// 	binary.LittleEndian.PutUint32(buf[i+1:], uint32(p.X))
	// 	binary.LittleEndian.PutUint32(buf[i+5:], uint32(p.Y))
	// 	binary.LittleEndian.PutUint16(buf[i+9:], uint16(keyboard))
	// 	i += 11
	// }
	return buf
}
