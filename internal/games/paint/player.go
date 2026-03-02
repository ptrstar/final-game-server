package paint

import "final-game-server/internal/shared"

type Player struct {
	Id       int
	Pos      *shared.Vec2
	PosVel   *shared.Vec2
	Angle    float32
	AngleVel float32
	Color    *shared.Col
	Keyboard map[byte]bool
}
