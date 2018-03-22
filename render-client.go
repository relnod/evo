package evo

import (
	"fmt"
	"log"

	"github.com/goxjs/glfw"
	"github.com/relnod/evo/num"
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
	window := system.NewWindow()
	renderer := system.NewRender(c.server.GetWorld())

	window.Init()
	renderer.Init()

	window.Window.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		if action == glfw.Release {
			x, y := w.GetCursorPos()
			log.Println(x, y)
			e := c.server.GetEntityAt(&num.Vec2{X: float32(x), Y: float32(y)})
			if e != nil {
				log.Printf("Entity:")
				fmt.Printf("Age: %f\n", e.Age)
				fmt.Println()
			}
		}
	})

	c.server.RegisterStream(func(w *world.World) {
		window.Update()
		renderer.UpdateWorld(w)
	})
}

// Start starts the client.
func (c *RenderClient) Start() {
	c.server.Start()
}
