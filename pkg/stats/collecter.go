package stats

import (
	"time"

	"github.com/relnod/evo/pkg/entity"
)

type EntityStatsSource interface {
	AnimalStats() *entity.DeathStats
	PlantStats() *entity.DeathStats
	ClearStats()
}

// IntervalCollecter collects stats in a given period.
type IntervalCollecter struct {
	interval int

	started time.Time
	stats   *Stats

	entityStatsSource EntityStatsSource
}

// NewIntervalCollector returns a new interval collector.
func NewIntervalCollector(entityStatsSource EntityStatsSource, seed int64, interval int) *IntervalCollecter {
	return &IntervalCollecter{
		interval:          interval,
		started:           time.Now(),
		stats:             NewStats(seed),
		entityStatsSource: entityStatsSource,
	}
}

// Update updates the stats if the interval period is over.
func (i *IntervalCollecter) Update(tick int, creatures []*entity.Creature) {
	if tick%i.interval != 0 {
		return
	}

	timeStat := newTimeStatFromCreatures(creatures)
	timeStat.Animal.DeathStats = *i.entityStatsSource.AnimalStats()
	timeStat.Plant.DeathStats = *i.entityStatsSource.PlantStats()

	i.stats.Running = time.Since(i.started) / (time.Millisecond * 1000)
	i.stats.Ticks = tick
	i.stats.Current = timeStat
	i.stats.OverTime.Add(timeStat)

	i.entityStatsSource.ClearStats()
}

// Stats returns the current stats.
func (i *IntervalCollecter) Stats() *Stats {
	return i.stats
}
