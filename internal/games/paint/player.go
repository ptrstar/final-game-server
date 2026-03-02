package paint

import "final-game-server/internal/shared"

type Player struct {
	Id       int           `json:"id"`
	Pos      *shared.Vec2  `json:"pos"`
	PosVel   *shared.Vec2  `json:"pos_vel"`
	Angle    float32       `json:"angle"`
	AngleVel float32       `json:"angle_vel"`
	Color    *shared.Col   `json:"color"`
	Keyboard map[byte]bool `json:"keyboard"`
}
