package evo

// Client defines an interface for a client
type Client interface {
	// Init initializes the client
	// TODO: should be able to return an error
	Init()

	// Start starts the client
	// TODO: should be able to return an error
	Start()
}
