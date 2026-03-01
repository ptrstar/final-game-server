package paint

import "final-game-server/internal/shared"

type PlayerInput struct {
	Data     []byte
	PlayerId player
}

type Player struct {
	Id       int
	Pos      *shared.Vec2
	Vel      *shared.Vec2
	Color    *shared.Col
	Keyboard map[byte]bool
}
