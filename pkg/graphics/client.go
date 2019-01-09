package graphics

import (
	"log"
	"time"

	"github.com/goxjs/glfw"
	"github.com/relnod/evo/pkg/entity"
	"github.com/relnod/evo/pkg/evo"
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

// Client implements evo.Consumer
type Client struct {
	producer evo.Producer

	creatures []*entity.Creature

	ticker *evo.Ticker

	window   *Window
	renderer *WorldRenderer
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
	width, height, err := c.producer.Size()
	if err != nil {
		log.Fatal(err)
		return
	}

	window := NewWindow(width, height)
	renderer := NewWorldRenderer()
	camera := NewCamera(width, height)

	window.Init()
	renderer.Init()

	camera.Connect(renderer)

	window.OnResize(func(width, height int) {
		renderer.SetSize(width, height)
		camera.SetSize(width, height)
		camera.Update()
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

	camera.Update()
	window.Update()

	c.producer.SubscribeEntitiesChanged(func(creatures []*entity.Creature) {
		c.creatures = creatures
	})

	c.window = window
	c.renderer = renderer
	c.ticker = evo.NewTicker(time.Second / 60)
}

// Start starts the client.
func (c *Client) Start() {
	go c.producer.Start()
	for tick := range c.ticker.C {
		_ = tick
		c.window.Update()
		c.renderer.Update(c.creatures)
	}
}

// Stop stops the graphics client.
func (c *Client) Stop() {
	c.producer.Stop()
	c.ticker.Stop()
}
