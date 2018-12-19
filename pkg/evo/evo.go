package evo

import (
	"github.com/google/uuid"

	"github.com/relnod/evo/pkg/world"
)

type WorldStream func(*world.World)

type Producer interface {
	Start() error
	GetWorld() (*world.World, error)
	SubscribeWorld(stream WorldStream) uuid.UUID
	UnsubscribeWorld(id uuid.UUID)
}

type Consumer interface {
	Init()
	Start() error
}
