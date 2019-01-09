package evo

import (
	"sync"
	"time"
)

// Ticker is a time ticker, with controlls for start, pause, and stop.
// It is safe to controll the ticker asynchron.
type Ticker struct {
	interval time.Duration

	C chan int

	// state of the ticker
	running bool
	pausing bool
	tick    int

	// m protects the state of the ticker.
	m *sync.Mutex
}

// NewTicker returns a new ticker.
func NewTicker(interval time.Duration) *Ticker {
	t := &Ticker{
		interval: interval,
		C:        make(chan int, 1),

		running: false,
		pausing: false,

		m: &sync.Mutex{},
	}
	go t.start()
	return t
}

// start starts the ticker.
func (t *Ticker) start() {
	t.running = true
	t.pausing = false

	for t.running {
		t.m.Lock()
		start := time.Now()
		if !t.pausing {
			t.tick++
			t.C <- t.tick
		}
		t.m.Unlock()

		time.Sleep(t.interval - time.Since(start))
	}
}

// Stop stops the ticker.
func (t *Ticker) Stop() { t.running = false; close(t.C) }

// Pause pauses the ticker.
func (t *Ticker) Pause() { t.pausing = true }

// Resume resumes the ticker.
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

// Interval returns the ticker delay.
func (t *Ticker) Interval() time.Duration {
	return t.interval
}

// SetInterval sets the Interval.
func (t *Ticker) SetInterval(interval time.Duration) {
	if interval < 0 {
		interval = 0
	}
	t.interval = interval
}
