package graphics

import (
	"log"

	"github.com/relnod/evo/pkg/evo"
	"github.com/relnod/evo/pkg/world"
)

// Client implements evo.Consumer
type Client struct {
	producer evo.Producer
}

// NewClient returns a new render client.
func NewClient(producer evo.Producer) *Client {
	return &Client{producer: producer}
}

// Init intitializes the window and renderer.
func (c *Client) Init() {
	w, err := c.producer.GetWorld()
	if err != nil {
		log.Fatal(err)
		return
	}

	window := NewWindow(int(w.Width), int(w.Height))
	renderer := NewWorldRenderer(w.Width, w.Height)
	window.OnResize(func(width, height int) {
		renderer.SetSize(width, height)
	})

	window.Init()
	renderer.Init()
	window.Update()

	c.producer.SubscribeWorld(func(w *world.World) {
		window.Update()
		renderer.Update(w)
	})
}

// Start starts the client.
func (c *Client) Start() {
	c.producer.Start()
}
