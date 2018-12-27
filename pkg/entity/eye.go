package entity

import (
	"math"
	"math/rand"

	"github.com/relnod/evo/pkg/math64"
)

type EyeDetection int

const (
	Biggest  EyeDetection = iota
	Smallest EyeDetection = iota
)

type Eye struct {
	Dir     math64.Vec2  `json:"dir"`
	Range   float64      `json:"range"`
	FOV     float64      `json:"fov"`
	Detects EyeDetection `json:"EyeDetection"`

	Count    int
	Detected float64
	Distance float64
}

func NewEye(eyeRange float64, detects EyeDetection) *Eye {
	return &Eye{
		Dir:     math64.Vec2{},
		Range:   eyeRange,
		FOV:     (80 / eyeRange * 40) * math.Pi / 180.0,
		Detects: detects,
	}
}

func NewRandomEye() *Eye {
	detects := Biggest
	if rand.Float64() > 0.5 {
		detects = Smallest
	}
	return NewEye(mutate(80.0, 1.0, 1.0), detects)
}

func (e *Eye) Sees(c *Creature) {
	e.Count++
	if e.Detects == Biggest {
		if c.Radius > e.Detected {
			e.Detected = c.Radius
		}
	} else if e.Detects == Smallest {
		if c.Radius < e.Detected {
			e.Detected = c.Radius
		}
	}
}

func (e *Eye) Reset() {
	e.Count = 0
	if e.Detects == Biggest {
		e.Detected = 0
	} else if e.Detects == Smallest {
		e.Detected = 9999999
	}
}
