package world

import (
	"github.com/relnod/evo/pkg/entity"
	"github.com/relnod/evo/pkg/math32"
	"github.com/relnod/evo/pkg/math32/collision"
)

// EdgeMode defines how the edge of the world is defined.
type EdgeMode int

const (
	// EdgeModeLoop defines the edge of a world as looping. This means that when
	// leaving on the left one enteres back on the right.
	EdgeModeLoop = 0
)

// StartMode defines how the world starts.
type StartMode int

const (
	// StartModeRandom generates world full of random initial entities.
	StartModeRandom = 0

	// StartModeFixed generates a world with one static entity.
	StartModeFixed = 1
)

// World holds all global world data
type World struct {
	Width  float32 `json:"width"`
	Height float32 `json:"height"`

	EdgeMode  EdgeMode  `json:"-"`
	StartMode StartMode `json:"-"`

	Creatures []*entity.Creature `json:"entities"`

	Static  []*entity.Creature `json:"-"`
	Dynamic []*entity.Creature `json:"-"`

	Cells       []*Cell `json:"-"`
	math32Cells int
	cellsPerRow int
	cellWidth   float32
	cellHeight  float32
}

// NewWorld returns a new world.
func NewWorld(width, height float32, edgeMode EdgeMode, startMode StartMode) *World {
	math32Cells := 36
	cellsPerRow := 6

	math32Entities := 1
	if startMode == StartModeRandom {
		math32Entities = 1000
	}

	w := &World{
		Width:  width,
		Height: height,

		EdgeMode:  edgeMode,
		StartMode: startMode,

		Creatures: make([]*entity.Creature, math32Entities),

		Static:  make([]*entity.Creature, 1),
		Dynamic: make([]*entity.Creature, 1),

		Cells:       make([]*Cell, math32Cells),
		math32Cells: math32Cells,
		cellsPerRow: cellsPerRow,
		cellWidth:   width / float32(math32Cells/cellsPerRow),
		cellHeight:  height / float32(math32Cells/cellsPerRow),
	}

	w.createCells()

	return w
}

func (w *World) createCells() {
	radius := (&math32.Vec2{X: w.cellWidth, Y: w.cellHeight}).Len()

	for row := 0; row < w.cellsPerRow; row++ {
		for col := 0; col < w.cellsPerRow; col++ {
			w.Cells[row*w.cellsPerRow+col] = &Cell{
				TopLeft: math32.Vec2{
					X: w.cellWidth * float32(row),
					Y: w.cellHeight * float32(col),
				},
				BotRight: math32.Vec2{
					X: w.cellWidth * float32(row+1),
					Y: w.cellHeight * float32(col+1),
				},

				Static:  make([]*entity.Creature, 0),
				Dynamic: make([]*entity.Creature, 0),
			}

			w.Cells[row*w.cellsPerRow+col].Center = math32.Vec2{
				X: w.Cells[row*w.cellsPerRow+col].TopLeft.X + w.cellWidth/2.0,
				Y: w.Cells[row*w.cellsPerRow+col].TopLeft.Y + w.cellHeight/2.0,
			}
			w.Cells[row*w.cellsPerRow+col].Radius = radius
		}
	}
}

// UpdateCells moves all creatures to it's corresponding cell.
func (w *World) UpdateCells() {
	for _, cell := range w.Cells {
		cell.Static = cell.Static[:0]
		cell.Dynamic = cell.Dynamic[:0]
	}

	for _, cell := range w.Cells {
		for _, creature := range w.Creatures {
			if collision.CircleCircle(&cell.Center, cell.Radius, &creature.Pos, creature.Radius) {
				if creature.Speed == 0 && creature.State != entity.StateChild {
					cell.Static = append(cell.Static, creature)
				} else {
					cell.Dynamic = append(cell.Dynamic, creature)
				}
			}
		}
	}
}

// FindCell returns the cell for the given position.
func (w *World) FindCell(pos *math32.Vec2) *Cell {
	x := pos.X / w.cellWidth
	y := pos.Y / w.cellHeight

	index := int(y)*w.cellsPerRow + int(x)

	if index > 0 && index > len(w.Cells)-1 {
		return nil
	}

	return w.Cells[index]
}

func (w *World) RemoveCreature(i int) {
	if i+1 >= len(w.Creatures) {
		w.Creatures = w.Creatures[:i]
		return
	}

	w.Creatures = append(w.Creatures[:i], w.Creatures[i+1:]...)
}
