package paint

import (
	"final-game-server/internal/engine"
	"final-game-server/internal/shared"
	"sync"
)

type PaintGame struct {
	Room         *engine.Room
	Join         chan *engine.Client
	Leave        chan *engine.Client
	Players      map[*engine.Client]*Player
	mu           sync.Mutex
	nextSGList   []*shared.ShareableGameState
	UpdateGSList chan []*shared.ShareableGameState
}

func NewPaintGame(Room *engine.Room) *PaintGame {
	game := &PaintGame{
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
