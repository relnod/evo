package evo

import (
	"time"

	"github.com/google/uuid"

	"github.com/relnod/evo/pkg/world"
)

// Stats describes runtime statistics of the simulation.
type Stats struct {
	start time.Time

	Seed    int64         `json:"seed"`
	Running time.Duration `json:"running"`

	Population int `json:"population"`
}

type WorldFn func(*world.World)

// Producer produces data.
type Producer interface {
	// Start start the producer.
	Start() error

	// Stop stops the producer.
	Stop() error

	// World reutnrs the current state of the world.
	World() (*world.World, error)

	// Stats returns some statistics of the world in its current state.
	Stats() (*Stats, error)

	// Ticks returns the ticks per second.
	// The producer should update after every tick.
	Ticks() (int, error)

	// SetTicks sets the ticks per second
	SetTicks(ticks int) error

	// SubscribeWorldChange subscribes to a world change.
	// Each time the world gets updated, the provided function gets called.
	// The returned unique id can be used to unsubscribe later.
	SubscribeWorldChange(fn WorldFn) uuid.UUID

	// UnsubscribeWorldChange ends a subscription to the world change.
	UnsubscribeWorldChange(id uuid.UUID)
}

// Consumer consumes data
type Consumer interface {
	Init()
	Start() error
	Stop() error
}
