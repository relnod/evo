package collision_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/relnod/evo/pkg/math64"
	"github.com/relnod/evo/pkg/math64/collision"
)

func TestCircleCircle(t *testing.T) {
	var tests = []struct {
		p1   *math64.Vec2
		r1   float64
		p2   *math64.Vec2
		r2   float64
		want bool
	}{
		{&math64.Vec2{X: 0, Y: 0}, 0.2, &math64.Vec2{X: 1, Y: 1}, 0.2, false},
		{&math64.Vec2{X: 0, Y: 0}, 1, &math64.Vec2{X: 1, Y: 1}, 1, true},
	}

	for i, test := range tests {
		got := collision.CircleCircle(test.p1, test.r1, test.p2, test.r2)
		assert.Equal(t, test.want, got, "Test case %d failed", i+1)
	}
}

func TestCirclePoint(t *testing.T) {
	var tests = []struct {
		p1   *math64.Vec2
		r1   float64
		p2   *math64.Vec2
		want bool
	}{
		{&math64.Vec2{X: 0, Y: 0}, 0.2, &math64.Vec2{X: 1, Y: 1}, false},
		{&math64.Vec2{X: 0, Y: 0}, 2, &math64.Vec2{X: 1, Y: 1}, true},
	}

	for i, test := range tests {
		got := collision.CirclePoint(test.p1, test.r1, test.p2)
		assert.Equal(t, test.want, got, "Test case %d failed", i+1)
	}
}
