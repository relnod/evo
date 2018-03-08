package system

import "github.com/relnod/evo/entity"

type System struct {
	Width  float32
	Height float32

	creatures []*entity.Creature
}

func NewSystem(width, height float32, numEntities int) *System {
	s := &System{
		Width:     width,
		Height:    height,
		creatures: make([]*entity.Creature, numEntities),
	}

	return s
}
