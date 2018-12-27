package collision_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/relnod/evo/pkg/math64"
	"github.com/relnod/evo/pkg/math64/collision"
)

func TestSquarePoint(t *testing.T) {
	var tests = []struct {
		topLeft     *math64.Vec2
		bottomRight *math64.Vec2
		p           *math64.Vec2
		want        bool
	}{
		{&math64.Vec2{X: 0, Y: 0}, &math64.Vec2{X: 1, Y: 1}, &math64.Vec2{X: 1.5, Y: 1}, false},
		{&math64.Vec2{X: 0, Y: 0}, &math64.Vec2{X: 1, Y: 1}, &math64.Vec2{X: 0.5, Y: 0.5}, true},
	}

	for i, test := range tests {
		got := collision.SquarePoint(test.topLeft, test.bottomRight, test.p)
		assert.Equal(t, test.want, got, "Test case %d failed", i+1)
	}
}
