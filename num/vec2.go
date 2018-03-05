package num

import (
	"math"
)

type Vec2 struct {
	X float32
	Y float32
}

func (v *Vec2) Len() float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y)))
}

func (v *Vec2) Norm() {
	l := v.Len()

	v.X /= l
	v.Y /= l
}
func (v *Vec2) Rotate(angle float64) {
	x := float64(v.X)
	y := float64(v.Y)

	v.X = float32(x*math.Cos(angle) - y*math.Sin(angle))
	v.Y = float32(x*math.Sin(angle) + y*math.Cos(angle))
}
