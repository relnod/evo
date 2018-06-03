package graphics

import (
	"github.com/relnod/evo"
	"github.com/relnod/evo/world"
)

// Client defines a render client.
type Client struct {
	server   evo.Server
	renderer Render
}

// NewClient returns a new render client.
func NewClient(server evo.Server) *Client {
	return &Client{server: server}
}

// Init intitializes the window and renderer.
func (c *Client) Init() {
	w := c.server.GetWorld()

	window := NewWindow(int(w.Width), int(w.Height))
	renderer := NewRender(w.Width, w.Height)

	window.Init()
	renderer.Init()

	c.server.RegisterStream(func(w *world.World) {
		window.Update()
		renderer.UpdateWorld(w)
	})
}

// Start starts the client.
func (c *Client) Start() {
	c.server.Start()
}
