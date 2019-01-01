package evo

import (
	"time"

	"github.com/relnod/evo/pkg/entity"
)

// Stats describes runtime statistics of the simulation.
type Stats struct {
	start time.Time

	Seed    int64         `json:"seed"`
	Running time.Duration `json:"running"`

	Current  *TimeStat        `json:"current"`
	OverTime *TimeStatHistory `json:"overtime"`

	Animal *entity.Stats `json:"animal"`
	Plant  *entity.Stats `json:"plant"`
}

func NewStats() *Stats {
	return &Stats{
		start: time.Now(),
		Seed:  time.Now().Unix(),
		Current: &TimeStat{
			Animal: &EntityTimeStat{},
			Plant:  &EntityTimeStat{},
		},
		OverTime: &TimeStatHistory{
			Population: make([]int, 0),
			Animal:     NewEntityTimeStatHistroy(),
			Plant:      NewEntityTimeStatHistroy(),
		},
	}

}

type TimeStat struct {
	Population int             `json:"population"`
	Animal     *EntityTimeStat `json:"animal"`
	Plant      *EntityTimeStat `json:"plant"`
}

type TimeStatHistory struct {
	Population []int                  `json:"population"`
	Animal     *EntityTimeStatHistory `json:"animal"`
	Plant      *EntityTimeStatHistory `json:"plant"`
}

func (t *TimeStatHistory) Add(stat *TimeStat) {
	t.Population = append(t.Population, stat.Population)
	t.Animal.Add(stat.Animal)
	t.Plant.Add(stat.Plant)
}

type EntityTimeStat struct {
	Population        int `json:"population"`
	HighestGeneration int `json:"highest_generation"`
}

func (e *EntityTimeStat) Add(c *entity.Creature) {
	e.Population++
	if e.HighestGeneration <= c.Consts.Generation {
		e.HighestGeneration = c.Consts.Generation
	}
}

type EntityTimeStatHistory struct {
	Population        []int `json:"population"`
	HighestGeneration []int `json:"highest_generation"`
}

func NewEntityTimeStatHistroy() *EntityTimeStatHistory {
	return &EntityTimeStatHistory{
		Population:        make([]int, 0),
		HighestGeneration: make([]int, 0),
	}
}

func (e *EntityTimeStatHistory) Add(stat *EntityTimeStat) {
	e.Population = append(e.Population, stat.Population)
	e.HighestGeneration = append(e.HighestGeneration, stat.HighestGeneration)
}
