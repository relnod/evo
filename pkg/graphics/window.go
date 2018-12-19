package graphics

import (
	"github.com/goxjs/gl"
	"github.com/goxjs/glfw"
)

type resizeCallback func(width, height int)

// Window wrappes a glfw.Window
type Window struct {
	Window *glfw.Window

	width  int
	height int

	resizeCallback resizeCallback
}

// NewWindow creates a new window
func NewWindow(width, height int) *Window {
	return &Window{width: width, height: height}
}

// Init inititializes the window.
func (w *Window) Init() {
	err := glfw.Init(gl.ContextWatcher)
	if err != nil {
		panic(err)
	}
	// defer glfw.Terminate()

	// width, height := glfw.GetPrimaryMonitor().GetPhysicalSize()
	window, err := glfw.CreateWindow(w.width, w.height, "Evo", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	width, height := window.GetSize()
	if w.resizeCallback != nil {
		w.resizeCallback(width, height)
	}

	glfw.SwapInterval(0)

	w.Window = window
}

// Update updates the window.
func (w *Window) Update() {
	w.Window.SwapBuffers()
	glfw.PollEvents()
}

func (w *Window) OnResize(cb resizeCallback) {
	w.resizeCallback = cb
}
