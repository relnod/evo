package world

import (
	"github.com/relnod/evo/pkg/entity"
	"github.com/relnod/evo/pkg/math32"
)

// Cell holds all entitties in a cell.
type Cell struct {
	TopLeft  math32.Vec2
	BotRight math32.Vec2

	Center math32.Vec2
	Radius float32

	Static  []*entity.Creature
	Dynamic []*entity.Creature
}
