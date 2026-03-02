package engine

import "final-game-server/internal/shared"

type GameItf interface {
	Run()
	AddPlayer(*Client)
	RemovePlayer(*Client)
	HandleInput(*ClientInput)
	Update()
	SerializeState() []byte
	GetShareableGameState() *shared.ShareableGameState
}
