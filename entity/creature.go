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

	Brain *deep.Neural

	Generation        int
	EnergyConsumption float32

	Alive        bool
	Saturation   float32
	SawFood      bool
	LastBread    float32
	Age          float32
	WantsToBreed bool
	Child        bool
}

func NewCreature(pos num.Vec2) *Creature {
	r := rand.Float32()*10 + 2.0

	brain := deep.NewNeural(&deep.Config{
		Inputs:     2,
		Layout:     []int{2, 2, 2},
		Activation: deep.ActivationSigmoid,
		Bias:       true,
		Weight:     deep.NewNormal(1.0, 0.0),
	})

	return newCreature(pos, r, brain, 0)

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

	r := e.Radius
	// r *= 1.0 + rand.Float32()/10.0

	return newCreature(e.Pos, r, deep.FromDump(dump), e.Generation+1)
}

func newCreature(pos num.Vec2, radius float32, brain *deep.Neural, generation int) *Creature {
	var speed float32 = 0.0
	if radius > 4.0 {
		speed = 5 / radius
	}

	return &Creature{
		Pos:    pos,
		Radius: radius,
		Dir:    randomDir(),
		Speed:  speed,

		Brain: brain,

		Generation:        generation,
		EnergyConsumption: rand.Float32() / 50,

		Alive:        true,
		Saturation:   5.0,
		SawFood:      false,
		LastBread:    -40,
		Age:          0,
		WantsToBreed: false,
		Child:        true,
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

	if e.Child {
		e.Pos.X += e.Dir.X
		e.Pos.Y += e.Dir.Y
		if e.Age > 0.5 {
			e.Child = false
		}
	} else {
		if e.Saturation > e.Radius*10 && (e.Age-e.LastBread) > 40 {
			e.WantsToBreed = true
		}

		if e.Speed > 0 {
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
		} else {
			e.Saturation += e.EnergyConsumption
		}
	}

	e.Age += 0.01

	e.SawFood = false
}

func (e *Creature) Collide(e2 *Creature) {
	if e.Radius > e2.Radius {
		e.Saturation += e2.Radius
		e2.Die()
	}
}

func (e *Creature) IsAlive() bool {
	return e.Alive
}

func (e *Creature) Die() {
	e.Alive = false
}
