package collision

import "github.com/relnod/evo/pkg/math32"

func SquarePoint(topLeft, botRight, p *math32.Vec2) bool {
	if p.X > topLeft.X && p.Y > topLeft.Y && p.X < botRight.X && p.Y < botRight.Y {
		return true
	}

	return false
}
