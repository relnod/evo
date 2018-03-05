package system

import "github.com/relnod/evo/entity"

type System struct {
	Width  float32
	Height float32

	creatures []*entity.Creature
	food      []*entity.Food
}

func NewSystem(width, height float32, numCreatures, numFood int) *System {
	s := &System{
		Width:     width,
		Height:    height,
		creatures: make([]*entity.Creature, numCreatures),
		food:      make([]*entity.Food, numFood),
	}

	return s
}
