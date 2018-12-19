package graphics

import (
	"log"

	"github.com/goxjs/glfw"
	"github.com/relnod/evo/pkg/evo"
	"github.com/relnod/evo/pkg/world"
)

type Renderer interface {
	UpdateViewport(zoom, x, y float32)
}

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
	camera := NewCamera(renderer)
	window.OnResize(func(width, height int) {
		renderer.SetSize(width, height)
	})
	window.OnKey(func(key glfw.Key) {
		switch key {
		case glfw.KeyKPAdd:
			camera.ZoomIn()
		case glfw.KeyKPSubtract:
			camera.ZoomOut()
		case glfw.KeyDown:
			camera.MoveDown()
		case glfw.KeyUp:
			camera.MoveUp()
		case glfw.KeyLeft:
			camera.MoveLeft()
		case glfw.KeyRight:
			camera.MoveRight()
		}
	})

	window.Init()
	renderer.Init()
	camera.Update()
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
