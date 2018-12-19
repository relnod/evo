package entity

import (
	"math/rand"

	deep "github.com/patrikeh/go-deep"
	"github.com/relnod/evo/pkg/config"
	"github.com/relnod/evo/pkg/math32"
)

type State int

const (
	StateChild State = iota
	StateAdult
	StateBreading
)

type Creature struct {
	Pos    math32.Vec2 `json:"pos"`
	Radius float32     `json:"radius"`
	Speed  float32     `json:"speed"`
	Dir    math32.Vec2 `json:"-"`

	Eye   *Eye         `json:"-"`
	Brain *deep.Neural `json:"-"`

	Alive     bool    `json:"-"`
	Energy    float32 `json:"-"`
	LastBread float32 `json:"-"`
	Age       float32 `json:"-"`
	State     State   `json:"-"`

	Consts Constants `json:"-"`
}

type Constants struct {
	Generation        int
	EnergyConsumption float32
	EnergyBreed       float32
	LifeExpectancy    float32
}

func NewCreature(pos math32.Vec2, radius float32) *Creature {
	return newCreature(pos, radius, newBrain(), 0, nil)
}

func (e *Creature) GetChild() *Creature {
	r := mutate(e.Radius, 0.1, 0.5)
	r = mutate(e.Radius, 0.9, 0.3)

	if r < 2.0 {
		r = 2.0
	} else if r > 10.0 {
		r = 10.0
	}

	return newCreature(e.Pos, r, e.Brain, e.Consts.Generation+1, e.Eye)
}

func newCreature(pos math32.Vec2, radius float32, brain *deep.Neural, generation int, eye *Eye) *Creature {
	var speed float32
	energyConsumption := rand.Float32() / 120
	energy := radius
	if radius > 4.0 {
		speed = mutate(15.0/(radius*radius), 0.2, 1.0)

		eyeRange := mutate(20.0, 0.2, 1.0)
		if eye != nil {
			eyeRange = mutate(eye.Range, 0.1, 0.2)
			eyeRange = mutate(eye.Range, 0.5, 0.1)
		}
		eye = NewEye(eyeRange)
		energyConsumption *= -1.0
		if brain == nil {
			brain = newBrain()
		}
	} else {
		brain = nil
	}

	if brain != nil {
		brain = newMutateBrain(brain)
	}

	return &Creature{
		Pos:    pos,
		Radius: radius,
		Dir:    randomDir(),
		Speed:  speed,

		Eye:   eye,
		Brain: brain,

		Alive:     true,
		Energy:    energy,
		LastBread: -30,
		Age:       0,
		State:     StateChild,

		Consts: Constants{
			Generation:        generation,
			EnergyConsumption: energyConsumption,
			EnergyBreed:       mutate(radius*radius*radius, 0.2, 0.9),
			LifeExpectancy:    mutate(100, 0.2, 1.0),
		},
	}
}

func newBrain() *deep.Neural {
	return deep.NewNeural(&deep.Config{
		Inputs:     2,
		Layout:     []int{2, 4, 4},
		Activation: deep.ActivationLinear,
		Bias:       true,
		Weight:     deep.NewNormal(1.0, 0.0),
	})
}

func newMutateBrain(brain *deep.Neural) *deep.Neural {
	dump := brain.Dump()
	for i := range dump.Weights {
		for j := range dump.Weights[i] {
			for k := range dump.Weights[i][j] {
				dump.Weights[i][j][k] = mutate64(dump.Weights[i][j][k], 0.1, 0.05)
				dump.Weights[i][j][k] = mutate64(dump.Weights[i][j][k], 0.5, 0.01)
			}
		}
	}

	return deep.FromDump(dump)
}

func mutate(val float32, fac float32, chance float32) float32 {
	if rand.Float32() > chance {
		return val
	}

	return val * (1.0 + (rand.Float32()-0.5)*fac)
}

func mutate64(val float64, fac float64, chance float64) float64 {
	if rand.Float64() > chance {
		return val
	}

	return val * (1.0 + (rand.Float64()-0.5)*fac)
}

func randomDir() math32.Vec2 {
	d := math32.Vec2{
		X: float32(rand.Float32()*2 - 1),
		Y: float32(rand.Float32()*2 - 1),
	}
	d.Norm()

	return d
}

func (e *Creature) Update() {
	if e.Energy <= 0 || e.Age > e.Consts.LifeExpectancy {
		e.Die()
	}

	if !e.Alive {
		return
	}

	switch e.State {
	case StateChild:
		e.Pos.X += e.Dir.X * config.WorldSpeed
		e.Pos.Y += e.Dir.Y * config.WorldSpeed

		if e.Age > 0.5 {
			e.State = StateAdult
		}
	case StateAdult:
		// if e.Energy > e.Consts.EnergyBreed && (e.Age-e.LastBread) > 40 {
		if e.Energy > e.Consts.EnergyBreed && (e.Age-e.LastBread) > 40 {
			e.State = StateBreading
		}

		if e.Speed > 0 {
			e.updateFromBrain()

			e.Pos.X += e.Dir.X * e.Speed * config.WorldSpeed
			e.Pos.Y += e.Dir.Y * e.Speed * config.WorldSpeed
		}

		e.Energy += e.Consts.EnergyConsumption * config.WorldSpeed
	}

	e.Age += 0.01 * config.WorldSpeed
}

func (e *Creature) updateFromBrain() {
	in1 := -0.9
	in2 := -0.9
	if e.Eye.Count > 0 {
		in1 += float64(e.Eye.Count) / 10.0
		if in1 > 0.9 {
			in1 = 0.9
		}

		if e.Eye.Biggest > e.Radius {
			in2 = 0.9
		}
		e.Eye.Reset()
	}

	out := e.Brain.Predict([]float64{in1, in2})
	if out[0] < 0 {
		rotation := 0.0
		if out[1] < -0.5 {
			rotation = 0.01
		} else if out[1] < 0 {
			rotation = 0.05
		} else if out[1] < 0.5 {
			rotation = 0.1
		} else {
			rotation = 0.14
		}

		if out[2] < 0 {
			rotation *= -1
		}

		e.Dir.Rotate(rotation)
		e.Dir.Norm()
	}

	if out[3] > 0 {
		e.Dir.X *= -1
		e.Dir.Y *= -1
	}

	e.Eye.Dir = e.Dir
}

func (e *Creature) Collide(e2 *Creature) {
	if e.Brain != nil && e2.Brain == nil {
		e.Energy += e2.Radius / 2.0 * 3.0
		e2.Die()
	} else if e.Radius > e2.Radius && !e.IsSameSpecies(e2) {
		e.Energy += e2.Radius / 2.0 * 3.0
		e2.Die()
	}
}

func (e *Creature) IsSameSpecies(e2 *Creature) bool {
	diff := e.Radius / e2.Radius
	return diff > 0.9 && diff < 1.1
}

func (e *Creature) IsAlive() bool {
	return e.Alive
}

func (e *Creature) Die() {
	e.Alive = false
}
