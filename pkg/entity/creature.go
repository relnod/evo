package entity

import (
	"math/rand"
	"time"

	deep "github.com/patrikeh/go-deep"

	"github.com/relnod/evo/pkg/config"
	"github.com/relnod/evo/pkg/math64"
)

// State defines the state of a creature.
type State int

// Defines all creature states.
const (
	StateChild State = iota
	StateAdult
	StateBreading
)

type Death int

const (
	DeathByAge    Death = 1
	DeathByHunger       = 3
	DeathByEaten        = 5
)

// Creature can either be moving (animal) or stand still (plant).
type Creature struct {
	// Current position in the world.
	Pos math64.Vec2 `json:"pos"`

	// The direction the creature is facing
	Dir    math64.Vec2 `json:"-"`
	Radius float64     `json:"radius"`
	Speed  float64     `json:"speed"`

	Eyes  []*Eye       `json:"eyes"`
	Brain *deep.Neural `json:"-"`

	Alive     bool    `json:"-"`
	Energy    float64 `json:"-"`
	LastBread float64 `json:"-"`
	Age       float64 `json:"-"`
	State     State   `json:"-"`

	Interactions int
	DeathBy      Death

	lastEaten time.Time

	Consts Constants `json:"-"`
}

type Constants struct {
	Generation        int
	EnergyConsumption float64
	EnergyBreed       float64
	LifeExpectancy    float64
}

func NewCreature(pos math64.Vec2, radius float64) *Creature {
	return newCreature(pos, radius, nil, 0, nil)
}

func (e *Creature) NewChild() *Creature {
	r := mutate(e.Radius, 0.1, 0.5)
	r = mutate(e.Radius, 1.5, 0.3)

	if r < 2.0 {
		r = 2.0
	} else if r > 10.0 {
		r = 10.0
	}

	return newCreature(e.Pos, r, e.Brain, e.Consts.Generation+1, e.Eyes)
}

func newCreature(pos math64.Vec2, radius float64, brain *deep.Neural, generation int, eyes []*Eye) *Creature {
	var speed float64
	var newEyes []*Eye
	// energyConsumption := mutate(rand.Float64()*radius, 0.1, 0.1)

	energyConsumption := (rand.NormFloat64()*0.1 + 1.0) / 300.0 * ((1.0*radius + 0.0*(radius*radius)) / 4)
	energy := radius

	// if radius > 4.0 {
	if brain != nil || radius > 2.0 && rand.Float64() > 0.99 {
		if brain == nil {
			generation = 0
		}
		speed = mutate(2/(radius), 0.2, 1.0)

		// If no eye exists create a new one.
		if len(eyes) == 0 {
			newEyes = append(newEyes, NewRandomEye())
		} else {
			// Mutate all existing eyes.
			for _, eye := range eyes {
				eyeRange := mutate(eye.Range, 0.5, 0.1)
				newEyes = append(newEyes, NewEye(eyeRange, eye.Detects))
			}

			// With a 10% chance a new eye appears.
			r := rand.Float64()
			if r > 0.9 {
				newEyes = append(newEyes, NewRandomEye())
			} else if r < 0.1 && len(newEyes) > 1 {
				newEyes = newEyes[:len(newEyes)-1]
			}
		}
		energyConsumption *= -1.0
		if brain == nil {
			brain = NewBrain(len(newEyes) * 2)
		} else {
			brain = NewMutatedBrain(brain, len(newEyes)*2)
		}

	} else {
		if brain != nil {
			generation = 0
		}
		brain = nil
	}

	return &Creature{
		Pos:    pos,
		Radius: radius,
		Dir:    randomDir(),
		Speed:  speed,

		Eyes:  newEyes,
		Brain: brain,

		Alive:     true,
		Energy:    energy,
		LastBread: -30,
		Age:       0,
		State:     StateChild,

		lastEaten: time.Now(),

		Consts: Constants{
			Generation:        generation,
			EnergyConsumption: energyConsumption,
			EnergyBreed:       mutate(math64.Poly(radius, 0, 0.5, 0.5), 0.05, 0.5),
			LifeExpectancy:    mutate(radius*radius*radius*radius, 0.2, 1.0),
		},
	}
}

func NewBrain(inputs int) *deep.Neural {
	return deep.NewNeural(&deep.Config{
		Inputs:     inputs,
		Layout:     []int{inputs, 4, 4},
		Activation: deep.ActivationLinear,
		Bias:       true,
		Weight:     deep.NewNormal(1.0, 0.0),
	})
}

func NewMutatedBrain(brain *deep.Neural, inputs int) *deep.Neural {
	min := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}
	newBrain := deep.NewNeural(&deep.Config{
		Inputs:     inputs,
		Layout:     []int{inputs, 4, 4},
		Activation: deep.ActivationLinear,
		Bias:       true,
		Weight:     deep.NewNormal(1.0, 0.0),
	}).Dump()
	dump := brain.Dump()
	for i := 0; i < min(len(dump.Weights), len(newBrain.Weights)); i++ {
		for j := 0; j < min(len(dump.Weights[i]), len(newBrain.Weights[i])); j++ {
			for k := 0; k < min(len(dump.Weights[i][j]), len(newBrain.Weights[i][j])); k++ {
				newBrain.Weights[i][j][k] = mutate64(dump.Weights[i][j][k], 0.1, 0.05)
				newBrain.Weights[i][j][k] = mutate64(dump.Weights[i][j][k], 0.5, 0.01)
			}
		}
	}

	return deep.FromDump(newBrain)
}

func mutate(val float64, fac float64, chance float64) float64 {
	if rand.Float64() > chance {
		return val
	}

	return val * (1.0 + (rand.Float64()-0.5)*fac)
}

func mutate64(val float64, fac float64, chance float64) float64 {
	if rand.Float64() > chance {
		return val
	}

	return val * (1.0 + (rand.Float64()-0.5)*fac)
}

func randomDir() math64.Vec2 {
	d := math64.Vec2{
		X: rand.Float64()*2 - 1,
		Y: rand.Float64()*2 - 1,
	}
	d.Norm()

	return d
}

// Update updates the state of the creature.
func (e *Creature) Update() {
	if !e.IsAlive() {
		return
	}

	if e.Energy <= 0 {
		e.Die(DeathByHunger)
		return
	}
	if e.Age > e.Consts.LifeExpectancy {
		e.Die(DeathByAge)
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
		if e.Energy > e.Consts.EnergyBreed && (e.Age-e.LastBread) > rand.NormFloat64()*0.2+(e.Consts.LifeExpectancy/3) {
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
	inputs := make([]float64, len(e.Eyes)*2)
	for i, eye := range e.Eyes {
		inputs[i*2] = -0.9
		inputs[i*2+1] = -0.9

		if eye.Count > 0 {
			inputs[i*2] += float64(eye.Count) / 10.0
			if inputs[i*2] > 0.9 {
				inputs[i*2] = 0.9
			}

			if eye.Detected > e.Radius {
				inputs[i*2+1] = 0.9
			}
		}
	}

	out := e.Brain.Predict(inputs)
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

	for _, eye := range e.Eyes {
		eye.Reset()
		eye.Dir = e.Dir
	}
}

// Collide gets called, when the creature collides with another creature.
func (e *Creature) Collide(e2 *Creature) {
	if time.Since(e.lastEaten) < time.Second {
		return
	}
	if e.Brain == nil {
		return
	}

	if e2.Brain == nil || (e.Radius > e2.Radius && !e.IsSameSpecies(e2)) {
		e.Interactions++
		e2.Interactions++
		e.Energy += e2.Radius * e2.Radius * e2.Radius * e2.Radius
		e2.Die(DeathByEaten)
		e.lastEaten = time.Now()
	}
}

// IsSameSpecies returns true if the difference between the radius of both
// creatures is less than 10%.
func (e *Creature) IsSameSpecies(e2 *Creature) bool {
	diff := e.Radius / e2.Radius
	return diff > 0.5 && diff < 1.5
}

// IsAlive returns true if the creature is alive.
func (e *Creature) IsAlive() bool {
	return e.Alive
}

// Die lets the creature die.
func (e *Creature) Die(death Death) {
	e.Alive = false
	e.DeathBy = death
}
