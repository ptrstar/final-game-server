package wasd

import (
	"encoding/binary"
	"sync"
)

type PlayerState struct {
	ID       uint8
	X, Y     int32
	keyboard map[byte]bool
}

type WASDGame struct {
	Players    map[interface{}]*PlayerState
	MaxPlayers int
	mu         sync.RWMutex
	nextID     uint8
}

func NewGame(max int) *WASDGame {
	return &WASDGame{
		Players:    make(map[interface{}]*PlayerState),
		MaxPlayers: max,
		nextID:     1,
	}
}

func (g *WASDGame) AddPlayer(c interface{}) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.Players[c] = &PlayerState{ID: g.nextID, X: 100, Y: 100, keyboard: make(map[byte]bool)}
	g.nextID++
}

func (g *WASDGame) RemovePlayer(c interface{}) {
	g.mu.Lock()
	defer g.mu.Unlock()
	delete(g.Players, c)
}

func (g *WASDGame) HandleInput(c interface{}, data []byte) {

	if len(data) == 0 {
		return
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	p, ok := g.Players[c]
	if !ok {
		return
	}

	key := data[0] >> 1
	var value bool = (data[0] & 0b01) != 0

	p.keyboard[key] = value
}

// 60Hz Gameupdates
func (g *WASDGame) Update() {
	for _, player := range g.Players {
		if player.keyboard[0] {
			player.Y -= 5
		}
		if player.keyboard[1] {
			player.X -= 5
		}
		if player.keyboard[2] {
			player.Y += 5
		}
		if player.keyboard[3] {
			player.X += 5
		}
	}
}

// Binary Format: [Count(1B)][ID(1B), X(4B), Y(4B), ...]
func (g *WASDGame) SerializeState() []byte {

	buf := make([]byte, 1+(len(g.Players)*11))
	buf[0] = byte(len(g.Players))

	i := 1
	for _, p := range g.Players {
		buf[i] = p.ID
		var keyboard uint16 = 0
		if p.keyboard[0] {
			keyboard |= 1
		}
		if p.keyboard[1] {
			keyboard |= 2
		}
		if p.keyboard[2] {
			keyboard |= 4
		}
		if p.keyboard[3] {
			keyboard |= 8
		}
		binary.LittleEndian.PutUint32(buf[i+1:], uint32(p.X))
		binary.LittleEndian.PutUint32(buf[i+5:], uint32(p.Y))
		binary.LittleEndian.PutUint16(buf[i+9:], uint16(keyboard))
		i += 11
	}
	return buf
}
