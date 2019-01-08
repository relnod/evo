package evo_test

import (
	"testing"

	"github.com/relnod/evo/pkg/evo"
)

// This Benchmark runs the simulation for 100 updates.
func BenchmarkSimulation(b *testing.B) {
	s := evo.NewSimulationFromSeed(1000, 1000, 1000, 2)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for i := 0; i < 100; i++ {
			s.Update()
		}
	}
}
