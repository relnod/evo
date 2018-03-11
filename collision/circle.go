package collision

import (
	"github.com/relnod/evo/num"
)

func CircleCircle(pos1 *num.Vec2, r1 float32, pos2 *num.Vec2, r2 float32) bool {
	return (&num.Vec2{X: pos1.X - pos2.X, Y: pos1.Y - pos2.Y}).Len() < r1+r2
}

func CirclePoint(pos1 *num.Vec2, r1 float32, pos2 *num.Vec2) bool {
	return (&num.Vec2{X: pos1.X - pos2.X, Y: pos1.Y - pos2.Y}).Len() < r1
}
