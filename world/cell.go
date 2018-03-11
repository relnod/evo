package world

import (
	"github.com/relnod/evo/entity"
	"github.com/relnod/evo/num"
)

// Cell holds all entitties in a cell.
type Cell struct {
	TopLeft  num.Vec2
	BotRight num.Vec2

	Center num.Vec2
	Radius float32

	Static  []*entity.Creature
	Dynamic []*entity.Creature
}
