package evo

import (
	"github.com/relnod/evo/system"
	"github.com/relnod/evo/world"
)

// RenderClient defines a render client.
type RenderClient struct {
	server   Server
	renderer system.Render
}

// NewRenderClient returns a new render client.
func NewRenderClient(server Server) *RenderClient {
	return &RenderClient{server: server}
}

// Init intitializes the window and renderer.
func (c *RenderClient) Init() {
	renderer := system.NewRender(c.server.GetWorld())
	renderer.Init()
	c.server.RegisterStream(func(w *world.World) {
		renderer.UpdateWorld(w)
	})
}

// Start starts the client.
func (c *RenderClient) Start() {
	c.server.Start()
}
