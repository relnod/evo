package math64

import "math"

// f(x) = ax + bx^2 + (cx^n) + c
func Poly(x float64, koefs ...float64) float64 {
	result := 0.0
	for i, k := range koefs {
		result += k * math.Pow(x, float64(i))
	}
	return result
}
