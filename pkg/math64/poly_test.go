package math64_test

import (
	"testing"

	"github.com/relnod/evo/pkg/math64"
	"github.com/stretchr/testify/assert"
)

func TestPoly(t *testing.T) {
	var tests = []struct {
		got  float64
		want float64
	}{
		{math64.Poly(1, 0, 0, 0), 0},
		{math64.Poly(2, 0, 1, 0), 2},
		{math64.Poly(2, 0, 0, 1), 4},
		{math64.Poly(2, 0, 0, 0, 1), 8},
		{math64.Poly(1, 1), 1},
		{math64.Poly(1, 1), 1},

		{math64.Poly(1, 0, 1), 1},
		{math64.Poly(1, 1, 1), 2},

		{math64.Poly(2, 3), 3},
	}

	for i, test := range tests {
		assert.Equal(t, test.want, test.got, "Test case %d failed", i+1)
	}
}
