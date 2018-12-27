package collision

import "github.com/relnod/evo/pkg/math64"

// CircleCircle perfomes a collision detection on two circles. Returnes true if
// they collide.
func CircleCircle(pos1 *math64.Vec2, r1 float64, pos2 *math64.Vec2, r2 float64) bool {
	return (&math64.Vec2{X: pos1.X - pos2.X, Y: pos1.Y - pos2.Y}).Len() < r1+r2
}

// CirclePoint checks if the given point is in the given circle.
func CirclePoint(pos1 *math64.Vec2, r float64, pos2 *math64.Vec2) bool {
	return (&math64.Vec2{X: pos1.X - pos2.X, Y: pos1.Y - pos2.Y}).Len() < r
}
