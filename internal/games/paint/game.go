package paint

import (
	"final-game-server/internal/engine"
	"final-game-server/internal/shared"
	"sync"
)

type PaintGame struct {
	Room    *engine.Room
	Join    chan *engine.Client
	Leave   chan *engine.Client
	Players map[*engine.Client]*Player
	mu      sync.Mutex
}

func NewPaintGame(Room *engine.Room) *PaintGame {
	game := &PaintGame{
		Room:    Room,
		Join:    make(chan *engine.Client),
		Leave:   make(chan *engine.Client),
		Players: make(map[*engine.Client]*Player),
		mu:      sync.Mutex{},
	}
	go game.Run()
	return game
}

func (g *PaintGame) Run() {

	for {
		select {
		case client := <-g.Join:
			g.mu.Lock()
			g.Players[client] = &Player{
				Id: client.Id,
			}
			g.mu.Unlock()
		}
	}

}

func (g *PaintGame) AddPlayer(client *engine.Client) {
	g.Join <- client
}

func (g *PaintGame) RemovePlayer(client *engine.Client) {

	g.mu.Lock()
	defer g.mu.Unlock()

	delete(g.Players, client)
}
func (g *PaintGame) HandleInput(client *engine.Client, data *engine.ClientInput) {

}
func (g *PaintGame) Update() {

}
func (g *PaintGame) SerializeState() []byte {
	g.mu.Lock()
	defer g.mu.Unlock()
	return []byte{}

}

func (g *PaintGame) GetShareableGameState() *shared.ShareableGameState {
	g.mu.Lock()
	defer g.mu.Unlock()

	return &shared.ShareableGameState{
		Type:        g.Room.Type,
		Id:          g.Room.Id,
		Capacity:    g.Room.Capacity,
		Status:      "Running",
		CanJoin:     true,
		PlayerCount: len(g.Players),
		CreatedAt:   g.Room.CreatedAt,
		Custom:      nil,
	}
}
