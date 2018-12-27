package math64

import (
	"math"
)

type Vec2 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func (v *Vec2) Len() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v *Vec2) Norm() {
	l := v.Len()

	v.X /= l
	v.Y /= l
}

func (v *Vec2) Rotate(angle float64) {
	x := v.X
	y := v.Y

	v.X = x*math.Cos(angle) - y*math.Sin(angle)
	v.Y = x*math.Sin(angle) + y*math.Cos(angle)
}

func Angle(v1, v2 *Vec2) float64 {
	return math.Cos((v1.X*v2.X + v1.Y*v2.Y) / (v1.Len() * v2.Len()))
}
