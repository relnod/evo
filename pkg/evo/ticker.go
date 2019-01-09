package evo

import (
	"sync"
	"time"
)

// updateFunc defines a function signature used for the  update callback.
type updateFunc func(tick int) error

// Ticker is a time ticker, that calls an update function for a certain ticks
// per second.
// It is safe to controll the ticker asynchron.
// TODO: improve interface of ticker.
type Ticker struct {
	ticksPerSecond int
	updateFunc     updateFunc

	// state of the ticker
	running bool
	pausing bool
	tick    int

	m *sync.Mutex
}

// NewTicker returns a new ticker.
func NewTicker(ticksPerSecond int, updateFunc updateFunc) *Ticker {
	return &Ticker{
		ticksPerSecond: ticksPerSecond,
		updateFunc:     updateFunc,

		running: false,
		pausing: false,

		m: &sync.Mutex{},
	}
}

// Start starts the ticker. This method only exists, if an error occurs during
// the updateFunc call or if ticker.stop() gets called.
// TODO: make start time more correct
func (t *Ticker) Start() error {
	t.running = true
	t.pausing = false

	for t.running {
		t.m.Lock()
		start := time.Now()
		if !t.pausing {
			t.tick++
			if err := t.updateFunc(t.tick); err != nil {
				return err
			}
		}
		t.m.Unlock()

		time.Sleep(time.Second/time.Duration(t.ticksPerSecond) - time.Since(start))
	}
	return nil

}

// Stop stops the ticker on the next tick.
func (t *Ticker) Stop() { t.running = false }

// Pause pauses the ticker on the next tick.
func (t *Ticker) Pause() { t.pausing = true }

// Resume resumes the ticker on the next tick.
func (t *Ticker) Resume() { t.pausing = false }

// TogglePauseResume toggles pause/resume.
func (t *Ticker) TogglePauseResume() {
	t.m.Lock()
	if t.pausing {
		t.Resume()
	} else {
		t.Pause()
	}
	t.m.Unlock()
}

// Lock locks the ticker.
func (t *Ticker) Lock() {
	t.m.Lock()
}

// Unlock unlocks the ticker.
func (t *Ticker) Unlock() {
	t.m.Unlock()
}

// SetUpdate sets the update callback.
func (t *Ticker) SetUpdate(updateFunc updateFunc) {
	t.updateFunc = updateFunc
}

// TicksPerSecond returns the ticks per second.
func (t *Ticker) TicksPerSecond() int {
	return t.ticksPerSecond
}

// SetTicksPerSecond sets the ticks per second.
func (t *Ticker) SetTicksPerSecond(ticksPerSecond int) {
	if ticksPerSecond <= 0 {
		ticksPerSecond = 1
	}
	t.ticksPerSecond = ticksPerSecond
}
