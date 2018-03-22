package graphics

import (
	"github.com/goxjs/gl"
	"github.com/goxjs/glfw"
)

// Window wrappes a glfw.Window
type Window struct {
	Window *glfw.Window
}

// NewWindow creates a new window
func NewWindow() *Window {
	return &Window{}
}

// Init inititializes the window.
func (w *Window) Init() {
	err := glfw.Init(gl.ContextWatcher)
	if err != nil {
		panic(err)
	}
	// defer glfw.Terminate()

	// width, height := glfw.GetPrimaryMonitor().GetPhysicalSize()
	window, err := glfw.CreateWindow(500, 500, "Evo", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	glfw.SwapInterval(0)

	w.Window = window
}

// Update updates the window.
func (w *Window) Update() {
	w.Window.SwapBuffers()
	glfw.PollEvents()
}
