package graphics

import (
	"log"

	"github.com/goxjs/glfw"
	"github.com/relnod/evo/pkg/evo"
	"github.com/relnod/evo/pkg/world"
)

type Renderer interface {
	UpdateViewport(zoom, x, y float64)
}

// Client implements evo.Consumer
type Client struct {
	producer evo.Producer

	window *Window
}

// NewClient returns a new render client.
func NewClient(producer evo.Producer) *Client {
	return &Client{producer: producer}
}

// Init intitializes the window and renderer.
func (c *Client) Init() {
	w, err := c.producer.World()
	if err != nil {
		log.Fatal(err)
		return
	}

	window := NewWindow(w.Width, w.Height)
	renderer := NewWorldRenderer(w.Width, w.Height)
	camera := NewCamera(renderer)
	window.OnResize(func(width, height int) {
		renderer.SetSize(width, height)
	})
	window.OnKey(func(key glfw.Key, mods glfw.ModifierKey) {
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
		case glfw.KeyW:
			if mods == glfw.ModControl {
				c.Stop()
			}
		}
	})

	window.Init()
	renderer.Init()
	camera.Update()
	window.Update()

	c.producer.SubscribeWorldChange(func(w *world.World) {
		window.Update()
		renderer.Update(w)
	})

	c.window = window
}

// Start starts the client.
func (c *Client) Start() {
	c.producer.Start()
}

// Stop stops the graphics client.
func (c *Client) Stop() {
	c.producer.Stop()
}
