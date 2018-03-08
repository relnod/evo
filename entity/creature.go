package entity

import (
	"math/rand"

	deep "github.com/patrikeh/go-deep"
	"github.com/relnod/evo/num"
)

type State int

const (
	StateChild State = iota
	StateAdult
	StateBreading
)

type Creature struct {
	Pos    num.Vec2
	Radius float32
	Speed  float32
	Dir    num.Vec2

	Eye   *Eye
	Brain *deep.Neural

	Generation        int
	EnergyConsumption float32

	Alive      bool
	Saturation float32
	LastBread  float32
	Age        float32
	State      State
}

func NewCreature(pos num.Vec2) *Creature {
	r := rand.Float32()*10 + 2.0

	return newCreature(pos, r, newBrain(), 0)
}

func newBrain() *deep.Neural {
	return deep.NewNeural(&deep.Config{
		Inputs:     2,
		Layout:     []int{2, 2, 2},
		Activation: deep.ActivationSigmoid,
		Bias:       true,
		Weight:     deep.NewNormal(1.0, 0.0),
	})
}

func newMutateBrain(brain *deep.Neural) *deep.Neural {
	dump := brain.Dump()
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

	return deep.FromDump(dump)
}

func (e *Creature) GetChild() *Creature {
	r := e.Radius * (1.0 + (rand.Float32()-0.5)/5.0)

	return newCreature(e.Pos, r, newMutateBrain(e.Brain), e.Generation+1)
}

func newCreature(pos num.Vec2, radius float32, brain *deep.Neural, generation int) *Creature {
	var speed float32 = 0.0
	var eye *Eye
	energyConsumption := rand.Float32() / 60
	if radius > 4.0 {
		speed = 5 / radius
		eye = &Eye{
			Dir:   num.Vec2{},
			Range: 20.0,
			FOV:   5.0,
		}
		energyConsumption *= -1.0
		if brain == nil {
			brain = newBrain()
		}
	}

	return &Creature{
		Pos:    pos,
		Radius: radius,
		Dir:    randomDir(),
		Speed:  speed,

		Eye:   eye,
		Brain: brain,

		Generation:        generation,
		EnergyConsumption: energyConsumption,

		Alive:      true,
		Saturation: 5.0,
		LastBread:  -40,
		Age:        0,
		State:      StateChild,
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

	switch e.State {
	case StateChild:
		e.Pos.X += e.Dir.X
		e.Pos.Y += e.Dir.Y

		if e.Age > 0.5 {
			e.State = StateAdult
		}
	case StateAdult:
		if e.Saturation > e.Radius*10 && (e.Age-e.LastBread) > 40 {
			e.State = StateBreading
		}

		if e.Speed > 0 {
			e.updateFromBrain()
		}

		e.Pos.X += e.Dir.X * e.Speed
		e.Pos.Y += e.Dir.Y * e.Speed

		e.Saturation += e.EnergyConsumption
	}

	e.Age += 0.01
}

func (e *Creature) updateFromBrain() {
	in1 := 0.0
	in2 := 0.0
	if e.Eye.Saw > 0 {
		in1 = 0.9
		in2 = float64(e.Eye.Saw) / 10.0
		e.Eye.Reset()
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

	e.Eye.Dir = e.Dir
}

func (e *Creature) Collide(e2 *Creature) {
	if e.Radius/e2.Radius > 1.1 {
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
