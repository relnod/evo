package entity

type DeathStats struct {
	Lifetime      uint32 `json:"lifetime"`
	Interactions  uint32 `json:"interactions"`
	Generation    uint32 `json:"generation"`
	DeathByAge    uint32 `json:"death_by_age"`
	DeathByHunger uint32 `json:"death_by_hunger"`
	DeathByEaten  uint32 `json:"death_by_eaten"`
}

func (d *DeathStats) Clear() {
	d.Lifetime = 0
	d.Interactions = 0
	d.Generation = 0
	d.DeathByEaten = 0
	d.DeathByAge = 0
	d.DeathByHunger = 0
}

func (d *DeathStats) Add(c *Creature) {
	d.Lifetime = addToAverage(d.Lifetime, int(c.Age))
	d.Interactions = addToAverage(d.Interactions, c.Interactions)
	d.Generation = addToAverage(d.Generation, c.Consts.Generation)

	switch c.DeathBy {
	case DeathByAge:
		d.DeathByAge++
	case DeathByEaten:
		d.DeathByEaten++
	case DeathByHunger:
		d.DeathByHunger++
	}
}

func addToAverage(average uint32, add int) uint32 {
	if average == 0 {
		return uint32(add)
	}
	return uint32(float64(average+uint32(add)) / 2.0)
}

type DeathStatsHistory struct {
	Lifetime      []uint32 `json:"death_lifetime"`
	Interactions  []uint32 `json:"death_interactions"`
	Generation    []uint32 `json:"death_generation"`
	DeathByAge    []uint32 `json:"death_by_age"`
	DeathByHunger []uint32 `json:"death_by_hunger"`
	DeathByEaten  []uint32 `json:"death_by_eaten"`
}

func NewDeathStatsHistory() *DeathStatsHistory {
	return &DeathStatsHistory{
		Lifetime:      make([]uint32, 0),
		Interactions:  make([]uint32, 0),
		Generation:    make([]uint32, 0),
		DeathByAge:    make([]uint32, 0),
		DeathByHunger: make([]uint32, 0),
		DeathByEaten:  make([]uint32, 0),
	}
}

func (d *DeathStatsHistory) Add(stat *DeathStats) {
	d.Lifetime = append(d.Lifetime, stat.Lifetime)
	d.Interactions = append(d.Interactions, stat.Interactions)
	d.Generation = append(d.Generation, stat.Generation)
	d.DeathByAge = append(d.DeathByAge, stat.DeathByAge)
	d.DeathByEaten = append(d.DeathByEaten, stat.DeathByEaten)
	d.DeathByHunger = append(d.DeathByHunger, stat.DeathByHunger)
}
