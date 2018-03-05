package entity

import (
	"math/rand"

	deep "github.com/patrikeh/go-deep"
	"github.com/relnod/evo/num"
)

type Creature struct {
	Pos    num.Vec2
	Radius float32
	Speed  float32
	Dir    num.Vec2

	Brain      *deep.Neural
	Alive      bool
	Saturation float32

	EnergyConsumption float32
	SawFood           bool
	LastBread         float32

	Generation int
	Age        float32
}

func NewCreature(pos num.Vec2) *Creature {
	return &Creature{
		Pos:               pos,
		Radius:            8.0,
		Dir:               randomDir(),
		Speed:             1.5,
		Alive:             true,
		Saturation:        5.0,
		EnergyConsumption: rand.Float32() / 50,
		Brain: deep.NewNeural(&deep.Config{
			Inputs:     2,
			Layout:     []int{2, 2, 2},
			Activation: deep.ActivationSigmoid,
			Bias:       true,
			Weight:     deep.NewNormal(1.0, 0.0),
		}),
		SawFood:    false,
		LastBread:  -30,
		Generation: 0,
		Age:        0,
	}
}

func (e *Creature) GetChild() *Creature {
	dump := e.Brain.Dump()
	for i := range dump.Weights {
		for j := range dump.Weights[i] {
			for k := range dump.Weights[i][j] {
				r := rand.Float32()
				if r < 0.05 {
					dump.Weights[i][j][k] += 0.1
				} else if r < 0.1 {
					dump.Weights[i][j][k] -= 0.1
				}
			}
		}
	}

	return &Creature{
		Pos:               e.Pos,
		Radius:            e.Radius,
		Dir:               randomDir(),
		Speed:             e.Speed,
		Alive:             true,
		Saturation:        5.0,
		EnergyConsumption: rand.Float32() / 50,
		Brain:             deep.FromDump(dump),
		SawFood:           false,
		LastBread:         -30,
		Generation:        e.Generation + 1,
		Age:               0,
	}
}

func randomDir() num.Vec2 {
	d := num.Vec2{
		X: float32(rand.Float32()*2 - 1),
		Y: float32(rand.Float32()*2 - 1),
	}
	d.Norm()

	return d
}

func (e *Creature) Update() {
	if e.Saturation <= 0 || e.Age > 100 {
		e.Die()
	}

	if !e.Alive {
		return
	}

	in1 := 0.1
	in2 := 0.9
	if e.SawFood {
		in1 = 0.9
		in2 = 0.1
	}
	out := e.Brain.Predict([]float64{in1, in2})
	if out[0] < 0.5 {
		if out[1] < 0.5 {
			e.Dir.Rotate(0.02)
		} else {
			e.Dir.Rotate(-0.02)
		}
		e.Dir.Norm()
	}

	e.Pos.X += e.Dir.X * e.Speed
	e.Pos.Y += e.Dir.Y * e.Speed

	e.Saturation -= e.EnergyConsumption
	e.Age += 0.01

	e.SawFood = false
}

func (e *Creature) IsAlive() bool {
	return e.Alive
}

func (e *Creature) Die() {
	e.Alive = false
}
