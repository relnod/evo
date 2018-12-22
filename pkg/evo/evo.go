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
}

type WorldFn func(*world.World)

// Producer produces data.
type Producer interface {
	// Start start the producer.
	Start() error

	// World reutnrs the current state of the world.
	World() (*world.World, error)

	// Stats returns the statistics of the world in its current state.
	Stats() (*Stats, error)

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
}
