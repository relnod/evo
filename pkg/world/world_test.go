package world_test

import (
	"testing"

	"github.com/relnod/evo/pkg/world"
	"github.com/stretchr/testify/assert"
)

func TestNewWorld(t *testing.T) {
	t.Run("cells get created correctly", func(tt *testing.T) {
		want := []*world.Cell{
			&world.Cell{},
		}
		w := world.NewWorld(10, 10)
		assert.Equal(tt, w.Cells, want)
	})
}

func TestFindCell(t *testing.T) {

}
