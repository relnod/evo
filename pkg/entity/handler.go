package entity

import (
	"math/rand"

	"github.com/relnod/evo/pkg/math64"
	"github.com/relnod/evo/pkg/math64/collision"
)

type grid map[int]map[int]*cell

func (g grid) clear() {
	for x := range g {
		for y := range g[x] {
			g[x][y].plants = make([]*Creature, 0)
		}
	}
}

func newGrid(width, height, gridSize int) grid {
	gridCols := height / gridSize
	gridRows := width / gridSize

	grid := make(map[int]map[int]*cell, gridCols)

	for x := 0; x < gridCols; x += gridSize {
		grid[x] = make(map[int]*cell, gridRows)
		for y := 0; y < gridRows; y += gridSize {
			grid[x][y] = &cell{}
		}
	}

	return grid
}

func (g grid) update(creatures []*Creature) {
	for _, c := range creatures {
		if c.Brain == nil {
			cell := g[int(c.Pos.X)][int(c.Pos.Y)]
			cell.plants = append(cell.plants, c)
		}
	}
}

type cell struct {
	plants []*Creature
}

// Handler implements the evo.EntityHandler.
type Handler struct {
	initialPopulation int
	width             int
	height            int

	grid grid

	animalStats *DeathStats
	plantStats  *DeathStats

	collectStats bool
}

// NewHandler returns a new entity handler.
func NewHandler(width, height int, initalPopulation int) *Handler {
	return &Handler{
		initialPopulation: initalPopulation,
		width:             width,
		height:            height,

		grid: newGrid(width, height, 10),

		animalStats:  &DeathStats{},
		plantStats:   &DeathStats{},
		collectStats: true,
	}
}

// InitPopulation initializes a population with a given count.
func (h *Handler) InitPopulation() []*Creature {
	var (
		creatures = make([]*Creature, h.initialPopulation)
	)

	for i := range creatures {
		radius := rand.Float64()*rand.Float64()*rand.Float64()*10 + 2.0

		creatures[i] = NewCreature(randomPosition(creatures, h.width, h.height, radius), radius)
	}

	return creatures
}

// UpdatePopulation updates all entities.
// Also adds new child entities and removes dead ones.
func (h *Handler) UpdatePopulation(creatures []*Creature) []*Creature {
	var remove []int
	for i, c := range creatures {
		c.Update()

		if !c.Alive {
			if h.collectStats {
				if c.Brain == nil {
					h.plantStats.Add(c)
				} else {
					h.animalStats.Add(c)
				}
			}
			remove = append(remove, i)
			continue
		}

		if c.State == StateBreading {
			c.State = StateAdult
			c.LastBread = c.Age
			c.Energy -= c.Radius
			for i := 0; i < rand.Intn(int(1/(c.Radius*c.Radius*c.Radius*c.Radius)*100)+1)+1; i++ {
				child := c.NewChild()
				if c.Energy-child.Energy > 0 {
					c.Energy -= child.Energy
					creatures = append(creatures, child)
				}
			}
		}
	}

	for _, i := range remove {
		creatures = RemoveEntity(creatures, i)
	}

	return creatures
}

func (h *Handler) AnimalStats() *DeathStats {
	return h.animalStats
}

func (h *Handler) PlantStats() *DeathStats {
	return h.plantStats
}

func (h *Handler) ClearStats() {
	h.plantStats.Clear()
	h.animalStats.Clear()
}

// RemoveEntity removes an entity at a given index.
func RemoveEntity(creatures []*Creature, i int) []*Creature {
	if i+1 >= len(creatures) {
		return creatures[:i]
	}

	return append(creatures[:i], creatures[i+1:]...)
}

func randomPosition(creatures []*Creature, width, height int, radius float64) math64.Vec2 {
	pos := math64.Vec2{
		X: rand.Float64()*(float64(width)-(2*radius)) + radius,
		Y: rand.Float64()*(float64(height)-(2*radius)) + radius,
	}

	for _, creature := range creatures {
		if creature == nil {
			continue
		}

		if collision.CircleCircle(&creature.Pos, creature.Radius, &pos, radius) {
			return randomPosition(creatures, width, height, radius)
		}
	}

	return pos
}
