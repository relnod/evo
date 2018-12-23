package world

import (
	"math"

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
	// StartModeRandom generates a world full of random entities.
	StartModeRandom = 0

	// StartModeFixed generates a world with one static entity
	StartModeFixed = 1
)

// Cell holds all entitties in a cell.
type Cell struct {
	TopLeft  math32.Vec2
	BotRight math32.Vec2

	Center math32.Vec2
	Radius float32

	Static  []*entity.Creature
	Dynamic []*entity.Creature
}

// World holds all global world data
type World struct {
	Width  float32 `json:"width"`
	Height float32 `json:"height"`

	Opts *Options

	Creatures []*entity.Creature `json:"entities"`

	Static  []*entity.Creature `json:"-"`
	Dynamic []*entity.Creature `json:"-"`

	Cells       []*Cell `json:"-"`
	numCells    int
	cellsPerRow int
	cellWidth   float32
	cellHeight  float32
}

// Options holds a set of optional options to configure the world.
type Options struct {
	EdgeMode EdgeMode

	StartMode StartMode

	EntitiesAtStart int
}

type IWorld interface {
	FindEntity() *entity.Creature
}

// NewWorld returns a new world.
func NewWorld(width, height float32) *World {
	return NewWorldWithOptions(width, height, &Options{})
}

func NewWorldWithOptions(width, height float32, opts *Options) *World {
	numCells := 36
	cellsPerRow := int(math.Sqrt(float64(numCells)))
	cellsPerRow = 6
	cellWidth := width / float32(numCells/cellsPerRow)
	cellHeight := height / float32(numCells/cellsPerRow)

	if opts.EntitiesAtStart == 0 {
		opts.EntitiesAtStart = 1
		if opts.StartMode == StartModeRandom {
			opts.EntitiesAtStart = 1000
		}
	}

	w := &World{
		Width:  width,
		Height: height,

		Opts: opts,

		Creatures: make([]*entity.Creature, opts.EntitiesAtStart),

		Static:  make([]*entity.Creature, 1),
		Dynamic: make([]*entity.Creature, 1),

		Cells:       CreateCells(cellWidth, cellHeight, cellsPerRow, numCells),
		numCells:    numCells,
		cellsPerRow: cellsPerRow,
		cellWidth:   cellWidth,
		cellHeight:  cellHeight,
	}

	return w
}

func CreateCells(cellWidth float32, cellHeight float32, cellsPerRow int, numCells int) []*Cell {
	radius := (&math32.Vec2{X: cellWidth, Y: cellHeight}).Len()

	cells := make([]*Cell, numCells)

	for row := 0; row < cellsPerRow; row++ {
		for col := 0; col < cellsPerRow; col++ {
			cells[row*cellsPerRow+col] = &Cell{
				TopLeft: math32.Vec2{
					X: cellWidth * float32(row),
					Y: cellHeight * float32(col),
				},
				BotRight: math32.Vec2{
					X: cellWidth * float32(row+1),
					Y: cellHeight * float32(col+1),
				},

				Static:  make([]*entity.Creature, 0),
				Dynamic: make([]*entity.Creature, 0),
			}

			cells[row*cellsPerRow+col].Center = math32.Vec2{
				X: cells[row*cellsPerRow+col].TopLeft.X + cellWidth/2.0,
				Y: cells[row*cellsPerRow+col].TopLeft.Y + cellHeight/2.0,
			}
			cells[row*cellsPerRow+col].Radius = radius
		}
	}

	return cells
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
func (w *World) findCell(pos *math32.Vec2) *Cell {
	x := pos.X / w.cellWidth
	y := pos.Y / w.cellHeight

	index := int(y)*w.cellsPerRow + int(x)

	if index > 0 && index > len(w.Cells)-1 {
		return nil
	}

	return w.Cells[index]
}

func (w *World) FindEntity(pos *math32.Vec2) *entity.Creature {
	cell := w.findCell(pos)

	if cell == nil {
		return nil
	}

	for _, c := range cell.Dynamic {
		if collision.CirclePoint(&c.Pos, c.Radius+5, pos) {
			return c
		}
	}

	for _, c := range cell.Static {
		if collision.CirclePoint(&c.Pos, c.Radius+5, pos) {
			return c
		}
	}

	return nil
}

func (w *World) RemoveEntity(i int) {
	if i+1 >= len(w.Creatures) {
		w.Creatures = w.Creatures[:i]
		return
	}

	w.Creatures = append(w.Creatures[:i], w.Creatures[i+1:]...)
}
