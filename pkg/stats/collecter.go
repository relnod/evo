package stats

import (
	"time"

	"github.com/relnod/evo/pkg/entity"
)

// IntervalCollecter collects stats in a given period.
type IntervalCollecter struct {
	interval int

	started time.Time
	stats   *Stats
}

// NewIntervalCollector returns a new interval collector.
func NewIntervalCollector(seed int64, interval int) *IntervalCollecter {
	return &IntervalCollecter{
		interval: interval,
		started:  time.Now(),
		stats:    NewStats(seed),
	}
}

// Update updates the stats if the interval period is over.
func (i *IntervalCollecter) Update(tick int, creatures []*entity.Creature) {
	if tick%i.interval != 0 {
		return
	}

	timeStat := newTimeStatFromCreatures(creatures)

	i.stats.Running = time.Since(i.started) / (time.Millisecond * 1000)
	i.stats.Ticks = tick
	i.stats.Current = timeStat
	i.stats.OverTime.Add(timeStat)
}

// Stats returns the current stats.
func (i *IntervalCollecter) Stats() *Stats {
	return i.stats
}
