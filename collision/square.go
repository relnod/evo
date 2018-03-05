package collision

import "github.com/relnod/evo/num"

func SquarePoint(topLeft, botRight, p *num.Vec2) bool {
	if p.X > topLeft.X && p.Y > topLeft.Y && p.X < botRight.X && p.Y < botRight.Y {
		return true
	}

	return false
}
