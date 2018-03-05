package collision_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/relnod/evo/collision"
	"github.com/relnod/evo/entity"
	"github.com/relnod/evo/num"
)

func BenchmarkCircleCircle(b *testing.B) {
	pos1 := &num.Vec2{X: 10.0, Y: 10.0}
	var r1 float32 = 5.0
	pos2 := &num.Vec2{Y: 10.0, X: 10.0}
	var r2 float32 = 5.0

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		collision.CircleCircle(pos1, r1, pos2, r2)
	}
}

func BenchmarkSystemCircleCircle(b *testing.B) {
	benchmarks := []struct {
		numSubs      int
		numCreatures int
	}{
		{4, 10000},
		{9, 10000},
		{16, 10000},
		{25, 10000},
		{36, 10000},
	}

	for _, bm := range benchmarks {
		b.Run(fmt.Sprintf("%d-%d", bm.numSubs, bm.numCreatures), func(b *testing.B) {
			benchmarkSystemCircleCircle(b, bm.numSubs, bm.numCreatures)
		})
	}
}

func benchmarkSystemCircleCircle(b *testing.B, numSubs, numCreatures int) {
	width := 500
	height := 500

	system := collision.NewSystem(float32(width), float32(height), numSubs)

	creatures := make([]*entity.Creature, numCreatures)
	for i := 0; i < numCreatures; i++ {
		creatures[i] = &entity.Creature{
			Pos: num.Vec2{
				X: float32(rand.Intn(width)),
				Y: float32(rand.Intn(height)),
			},
			Radius: 5.0,
		}
		system.AddCreature(creatures[i])
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		system.Update()
	}
}
