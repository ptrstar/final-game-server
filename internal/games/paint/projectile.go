package paint

import (
	"final-game-server/internal/engine"
	"final-game-server/internal/shared"
)

type Projectile struct {
	Pos   *shared.Vec2
	Vel   *shared.Vec2
	Col   *shared.Col
	TTL   float32
	Owner *engine.Client
}
