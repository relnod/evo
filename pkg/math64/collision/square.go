package collision

import "github.com/relnod/evo/pkg/math64"

// SquarePoint checks if the given point is inside the given square.
func SquarePoint(topLeft, botRight, p *math64.Vec2) bool {
	if p.X > topLeft.X && p.Y > topLeft.Y && p.X < botRight.X && p.Y < botRight.Y {
		return true
	}

	return false
}
