package evo

import (
	"sync"
	"time"
)

// updateFunc defines a function signature used for the  update callback.
type updateFunc func(tick int) error

// ticker is a time ticker, that calls an update function for a certain ticks
// per second.
// It is safe to controll the ticker asynchron.
// It is also possible to set a function callback, that will always be called,
// while running. Even when the ticker is paused.
type ticker struct {
	ticksPerSecond   int
	updateFunc       updateFunc
	alwaysUpdateFunc updateFunc

	// state of the ticker
	running bool
	pausing bool
	tick    int

	m *sync.Mutex
}

func newTicker(ticksPerSecond int, updateFunc updateFunc) *ticker {
	return &ticker{
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
func (t *ticker) Start() error {
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
		if err := t.alwaysUpdateFunc(t.tick); err != nil {
			return err
		}
		t.m.Unlock()

		time.Sleep(time.Second/time.Duration(t.ticksPerSecond) - time.Since(start))
	}
	return nil

}

// Stop stops the ticker on the next tick.
func (t *ticker) Stop() { t.running = false }

// Pause pauses the ticker on the next tick.
func (t *ticker) Pause() { t.pausing = true }

// Resume resumes the ticker on the next tick.
func (t *ticker) Resume() { t.pausing = true }

// TogglePauseResume toggles pause/resume.
func (t *ticker) TogglePauseResume() {
	t.m.Lock()
	if t.pausing {
		t.Resume()
	} else {
		t.Pause()
	}
	t.m.Unlock()
}

// Lock locks the ticker.
func (t *ticker) Lock() {
	t.m.Lock()
}

// Unlock unlocks the ticker.
func (t *ticker) Unlock() {
	t.m.Unlock()
}

// SetUpdate sets the update callback.
func (t *ticker) SetUpdate(updateFunc updateFunc) {
	t.updateFunc = updateFunc
}

// SetUpdate sets update callback that always gets called.
func (t *ticker) SetAlwaysUpdate(updateFunc updateFunc) {
	t.alwaysUpdateFunc = updateFunc
}

// TicksPerSecond returns the ticks per second.
func (t *ticker) TicksPerSecond() int {
	return t.ticksPerSecond
}

// SetTicksPerSecond sets the ticks per second.
func (t *ticker) SetTicksPerSecond(ticksPerSecond int) {
	if ticksPerSecond <= 0 {
		ticksPerSecond = 1
	}
	t.ticksPerSecond = ticksPerSecond
}
