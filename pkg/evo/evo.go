package evo

import (
	"github.com/google/uuid"

	"github.com/relnod/evo/pkg/entity"
)

// EntitiesChangedFn defnies a callback function for entities.
type EntitiesChangedFn func([]*entity.Creature)

// Producer produces data.
type Producer interface {
	// Start start the producer.
	Start() error

	// Stop stops the producer.
	Stop() error

	// PauseResume Toggles pause/resume.
	PauseResume() error

	// Restart restarts the simulation.
	Restart() error

	// Size returns the size of the simulation.
	Size() (width int, height int, err error)

	// Creatures returns all creatures in their current state.
	Creatures() ([]*entity.Creature, error)

	// Stats returns some statistics of the world in its current state.
	Stats() (*Stats, error)

	// Ticks returns the ticks per second.
	// The producer should update after every tick.
	Ticks() (int, error)

	// SetTicks sets the ticks per second.
	SetTicks(ticks int) error

	// SubscribeEntitiesChanged subscribes to changes of entities.
	// Each time the entities get updated, the provided function gets called.
	// The returned unique id can be used to unsubscribe later.
	SubscribeEntitiesChanged(fn EntitiesChangedFn) uuid.UUID

	// UnsubscribeWorldChange ends a subscription to the world change.
	UnsubscribeEntitiesChanged(id uuid.UUID)
}

// Consumer consumes data
type Consumer interface {
	Init()
	Start() error
	Stop() error
}
