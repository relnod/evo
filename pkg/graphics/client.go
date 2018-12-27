package graphics

import (
	"log"

	"github.com/goxjs/glfw"
	"github.com/relnod/evo/pkg/evo"
	"github.com/relnod/evo/pkg/world"
)

const usage = `Keybindings:
Ctrl-w        Stop client
Ctrl-q        Stop client

r             Restarts the remote simulation
Space         Toggles pause/resume

Ctrl-Add      Increase ticks per seconds
Ctrl-Subtract Decrease ticks per seconds

ArrwoLeft     Move camera left
ArrwoRight    Move camera right
ArrwoUp       Move camera up
ArrwoDown     Move camera down
Add           Zoom in
Subtract      Zoom out`

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

// Usage returns the usage for the graphics client.
func (c *Client) Usage() string {
	return usage
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
			if mods == glfw.ModControl {
				ticks, err := c.producer.Ticks()
				if err != nil {
					log.Fatal(err.Error())
				}
				c.producer.SetTicks(ticks + 10)
			} else {
				camera.ZoomIn()
			}
		case glfw.KeyKPSubtract:
			if mods == glfw.ModControl {
				ticks, err := c.producer.Ticks()
				if err != nil {
					log.Fatal(err.Error())
				}
				c.producer.SetTicks(ticks - (ticks / 2))
			} else {
				camera.ZoomOut()
			}
		case glfw.KeyDown:
			camera.MoveDown()
		case glfw.KeyUp:
			camera.MoveUp()
		case glfw.KeyLeft:
			camera.MoveLeft()
		case glfw.KeyRight:
			camera.MoveRight()
		case glfw.KeyQ:
			if mods == glfw.ModControl {
				c.Stop()
			}
		case glfw.KeyW:
			if mods == glfw.ModControl {
				c.Stop()
			}
		case glfw.KeyR:
			c.producer.Restart()
		case glfw.KeySpace:
			c.producer.PauseResume()
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
