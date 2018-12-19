package evo

import (
	"github.com/google/uuid"

	"github.com/relnod/evo/pkg/world"
)

// WorldStream defines a function that updates the world model.
type WorldStream func(*world.World)

// Producer produces data.
type Producer interface {
	Start() error
	GetWorld() (*world.World, error)
	SubscribeWorld(stream WorldStream) uuid.UUID
	UnsubscribeWorld(id uuid.UUID)
}

// Consumer consumes data
type Consumer interface {
	Init()
	Start() error
}
