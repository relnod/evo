package collision

import "github.com/relnod/evo/pkg/math32"

func CircleCircle(pos1 *math32.Vec2, r1 float32, pos2 *math32.Vec2, r2 float32) bool {
	return (&math32.Vec2{X: pos1.X - pos2.X, Y: pos1.Y - pos2.Y}).Len() < r1+r2
}

func CirclePoint(pos1 *math32.Vec2, r float32, pos2 *math32.Vec2) bool {
	return (&math32.Vec2{X: pos1.X - pos2.X, Y: pos1.Y - pos2.Y}).Len() < r
}
