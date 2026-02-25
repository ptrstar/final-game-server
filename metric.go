package main

type Metric struct {
	players int
	tps     float32
	pIdle   float32
}

// how to do 'inheritance'

type GameItf interface {
	Start()
	GetScore() int
}

type BaseGame struct {
	ID    string
	Score int
}

func (b *BaseGame) GetScore() int { return b.Score }

type Chess struct {
	BaseGame // Embedding: chess now has all fields and methods of BaseGame
}

func (p *Chess) Start() {
	// defines chess specific start logic
}

func Test() {
	games := make(map[string]GameItf)

	games["table1"] = &Chess{BaseGame: BaseGame{ID: "c1", Score: 0}}

	// You can now call common methods on any item in the map
	for _, g := range games {
		g.Start()
	}
}
