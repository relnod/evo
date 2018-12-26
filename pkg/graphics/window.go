package graphics

import (
	"github.com/goxjs/gl"
	"github.com/goxjs/glfw"
)

type resizeCallback func(width, height int)
type keyCallback func(key glfw.Key, mods glfw.ModifierKey)

// Window wrappes a glfw.Window
type Window struct {
	window *glfw.Window

	width  int
	height int

	resizeCallback resizeCallback
	keyCallback    keyCallback
}

// NewWindow creates a new window
func NewWindow(width, height int) *Window {
	return &Window{
		width:  width,
		height: height,
	}
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
	window.SetKeyCallback(func(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if action != glfw.Press && action != glfw.Repeat {
			return
		}

		if w.keyCallback != nil {
			w.keyCallback(key, mods)
		}
	})

	w.window = window
}

// Update updates the window.
func (w *Window) Update() {
	w.window.SwapBuffers()
	glfw.PollEvents()
}

// Close closes the window.
func (w *Window) Close() {
	w.window.Destroy()
}

func (w *Window) OnResize(cb resizeCallback) {
	w.resizeCallback = cb
}

func (w *Window) OnKey(cb keyCallback) {
	w.keyCallback = cb
}
