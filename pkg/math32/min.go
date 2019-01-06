package math32

// Min returns the smaller of two numbers.
func Min(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}
