package entity_test

import (
	"testing"

	"github.com/relnod/evo/pkg/entity"
	"github.com/stretchr/testify/assert"
)

func TestFindOldest(t *testing.T) {
	plant := &entity.Creature{Consts: entity.Constants{Generation: 100}}
	animalAge1 := &entity.Creature{Brain: entity.NewBrain(1), Consts: entity.Constants{Generation: 1}}
	animalAge10 := &entity.Creature{Brain: entity.NewBrain(1), Consts: entity.Constants{Generation: 10}}
	tests := []struct {
		population []*entity.Creature
		want       *entity.Creature
	}{
		{
			[]*entity.Creature{}, nil,
		},
		{
			[]*entity.Creature{plant}, nil,
		},
		{
			[]*entity.Creature{plant, animalAge1, animalAge10}, animalAge10,
		},
	}

	for _, test := range tests {
		got := entity.FindOldest(test.population)
		assert.Equal(t, test.want, got)
	}
}
