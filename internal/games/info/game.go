package info

import (
	"encoding/json"
	"final-game-server/internal/engine"
	"final-game-server/internal/shared"
	"log"
	"sync"
)

type Info struct {
	Room         *engine.Room
	Join         chan *engine.Client
	Leave        chan *engine.Client
	Players      map[*engine.Client]*Player
	mu           sync.Mutex
	nextSGList   []*shared.ShareableGameState
	UpdateGSList chan []*shared.ShareableGameState
}

func NewInfo(Room *engine.Room) *Info {
	game := &Info{
		Room:         Room,
		Join:         make(chan *engine.Client),
		Leave:        make(chan *engine.Client),
		Players:      make(map[*engine.Client]*Player),
		mu:           sync.Mutex{},
		nextSGList:   make([]*shared.ShareableGameState, 0),
		UpdateGSList: make(chan []*shared.ShareableGameState),
	}
	go game.Run()
	return game
}

func (g *Info) Run() {

	for {
		select {
		case NewGSList := <-g.UpdateGSList:
			g.mu.Lock()
			g.nextSGList = NewGSList
			g.mu.Unlock()

		case client := <-g.Join:
			g.mu.Lock()
			g.Players[client] = &Player{
				Id: client.Id,
			}
			g.mu.Unlock()
		}
	}

}

func (g *Info) AddPlayer(client *engine.Client) {
	g.Join <- client
}

func (g *Info) RemovePlayer(client *engine.Client) {

	g.mu.Lock()
	defer g.mu.Unlock()

	delete(g.Players, client)
}
func (g *Info) HandleInput(client *engine.Client, data *engine.ClientInput) {

}
func (g *Info) Update() {

}
func (g *Info) SerializeState() []byte {

	g.mu.Lock()
	defer g.mu.Unlock()

	payload, err := json.Marshal(g.nextSGList)
	if err != nil {
		log.Printf("State serialization failed")
		return []byte{}
	}
	return payload
}

func (g *Info) GetShareableGameState() *shared.ShareableGameState {
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

func (g *Info) SetGSList(GSList []*shared.ShareableGameState) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.UpdateGSList <- GSList
}
