package evo

import (
	"github.com/relnod/evo/world"
	uuid "github.com/satori/go.uuid"
)

// Server defines an interface for a server.
type Server interface {
	// Start starts the server.
	Start()

	// GetWorld returns the current world object.
	GetWorld() *world.World

	// RegisterStream enables to register for a world stream.
	RegisterStream(stream Stream) uuid.UUID

	// UnRegisterStream removes the stream from the server.
	UnRegisterStream(id uuid.UUID)
}
