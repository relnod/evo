package world

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/relnod/evo/interal/testutil"
	"github.com/relnod/evo/pkg/entity"
	"github.com/relnod/evo/pkg/math64"
	"github.com/relnod/evo/pkg/math64/collision"
)

func testCollisionDetector(t *testing.T, collisionDetector CollisionDetector) {
	c1 := &entity.Creature{Speed: 1, Radius: 2, Pos: math64.Vec2{X: 1, Y: 1}}
	c2 := &entity.Creature{Speed: 1, Radius: 2, Pos: math64.Vec2{X: 1, Y: 2}}

	cLeft := &entity.Creature{Speed: 1, Radius: 1, Pos: math64.Vec2{X: -1, Y: 5}}
	cRight := &entity.Creature{Speed: 1, Radius: 1, Pos: math64.Vec2{X: 11, Y: 5}}
	cTop := &entity.Creature{Speed: 1, Radius: 1, Pos: math64.Vec2{X: 5, Y: -1}}
	cBot := &entity.Creature{Speed: 1, Radius: 1, Pos: math64.Vec2{X: 5, Y: 11}}
	cOutOfBoundsChild := &entity.Creature{State: entity.StateChild, Radius: 1, Pos: math64.Vec2{X: 100, Y: 100}}

	eye := &entity.Eye{Range: 2, FOV: math.Pi}
	cEye := &entity.Creature{Speed: 1, Dir: math64.Vec2{X: 1, Y: 1}, Pos: math64.Vec2{X: 1, Y: 1}, Eyes: []*entity.Eye{eye}}
	cSeen := &entity.Creature{Radius: 1, Pos: math64.Vec2{X: 2, Y: 2}}
	cNotSeen1 := &entity.Creature{Radius: 1, Pos: math64.Vec2{X: 0, Y: 0}}
	cNotSeen2 := &entity.Creature{Radius: 1, Pos: math64.Vec2{X: 5, Y: 5}}

	tests := []struct {
		desc       string
		population []*entity.Creature
		want       []Collision
	}{
		{
			"no collisions with empty population - who would have thought",
			[]*entity.Creature{},
			nil,
		},
		{
			"detects collision between a creature and the world border",
			[]*entity.Creature{cLeft, cRight, cTop, cBot, cOutOfBoundsChild},
			[]Collision{
				&creatureBorderCollision{cLeft, collision.LEFT, 10, 10},
				&creatureBorderCollision{cRight, collision.RIGHT, 10, 10},
				&creatureBorderCollision{cTop, collision.TOP, 10, 10},
				&creatureBorderCollision{cBot, collision.BOT, 10, 10},
				&creatureBorderCollision{cOutOfBoundsChild, collision.RIGHT, 10, 10},
			},
		},
		{
			"detects collision between two creatures, if they are moving",
			[]*entity.Creature{c1, c2},
			[]Collision{
				&creatureCreatureCollision{c1, c2},
				&creatureCreatureCollision{c2, c1},
			},
		},
		{
			"detects collision between eye fov and a creature",
			[]*entity.Creature{cEye, cSeen, cNotSeen1, cNotSeen2},
			[]Collision{
				&eyeCreatureCollision{eye, cSeen},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			got := collisionDetector.DetectCollisions(test.population)
			assert.Equal(tt, test.want, got)
		})
	}
}

func benchmarkCollisionDetector(b *testing.B, collisionDetector CollisionDetector) {
	population := testutil.Population(1000)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		collisionDetector.DetectCollisions(population)
	}
}

func TestSimpleCollisionDetector(t *testing.T) {
	testCollisionDetector(t, NewSimpleCollisionDetector(10, 10))
}

func BenchmarkSimpleCollisionDetector(b *testing.B) {
	benchmarkCollisionDetector(b, NewSimpleCollisionDetector(10, 10))
}
