package stats

import (
	"time"

	"github.com/relnod/evo/pkg/entity"
)

// Stats describes runtime statistics of the simulation.
type Stats struct {
	Seed    int64         `json:"seed"`
	Running time.Duration `json:"running"`
	Ticks   int           `json:"ticks"`

	Current  *timeStat        `json:"current"`
	OverTime *timeStatHistory `json:"overtime"`
}

// NewStats returns a new stats object.
func NewStats(seed int64) *Stats {
	return &Stats{
		Seed: seed,
		Current: &timeStat{
			Animal: &entityTimeStat{},
			Plant:  &entityTimeStat{},
		},
		OverTime: &timeStatHistory{
			Population: make([]int, 0),
			Animal:     newEntityTimeStatHistroy(),
			Plant:      newEntityTimeStatHistroy(),
		},
	}

}

type timeStat struct {
	Population int             `json:"population"`
	Animal     *entityTimeStat `json:"animal"`
	Plant      *entityTimeStat `json:"plant"`
}

func newTimeStatFromCreatures(creatures []*entity.Creature) *timeStat {
	t := &timeStat{
		Population: len(creatures),
		Animal:     &entityTimeStat{},
		Plant:      &entityTimeStat{},
	}

	for _, c := range creatures {
		if c.Brain == nil {
			t.Plant.Add(c)
		} else {
			t.Animal.Add(c)
		}

	}
	return t
}

type timeStatHistory struct {
	Population []int                  `json:"population"`
	Animal     *entityTimeStatHistory `json:"animal"`
	Plant      *entityTimeStatHistory `json:"plant"`
}

func (t *timeStatHistory) Add(stat *timeStat) {
	t.Population = append(t.Population, stat.Population)
	t.Animal.Add(stat.Animal)
	t.Plant.Add(stat.Plant)
}

type entityTimeStat struct {
	Population        int `json:"population"`
	HighestGeneration int `json:"highest_generation"`
}

func (e *entityTimeStat) Add(c *entity.Creature) {
	e.Population++
	if e.HighestGeneration <= c.Consts.Generation {
		e.HighestGeneration = c.Consts.Generation
	}
}

type entityTimeStatHistory struct {
	Population        []int `json:"population"`
	HighestGeneration []int `json:"highest_generation"`
}

func newEntityTimeStatHistroy() *entityTimeStatHistory {
	return &entityTimeStatHistory{
		Population:        make([]int, 0),
		HighestGeneration: make([]int, 0),
	}
}

func (e *entityTimeStatHistory) Add(stat *entityTimeStat) {
	e.Population = append(e.Population, stat.Population)
	e.HighestGeneration = append(e.HighestGeneration, stat.HighestGeneration)
}
