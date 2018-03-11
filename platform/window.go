package platform

import (
	"strconv"
	"time"

	"github.com/gopherjs/gopherjs/js"
	"github.com/relnod/evo/num"
)

type Window struct {
	Canvas *js.Object
	Width  int
	Height int

	keyListeners   []KeyListener
	mouseListeners []MouseListener
}

func NewWindow() *Window {
	document := js.Global.Get("document")

	// @todo: instead wait till dom is ready
	time.Sleep(time.Millisecond)

	window := js.Global.Get("window")
	width, err := strconv.Atoi(window.Get("innerWidth").String())
	if err != nil {
		// @todo
	}
	height, err := strconv.Atoi(window.Get("innerHeight").String())
	if err != nil {
		// @todo
	}

	canvas := document.Call("createElement", "canvas")
	canvas.Set("width", width)
	canvas.Set("height", height)

	body := document.Get("body")
	body.Get("style").Set("margin", 0)
	body.Call("appendChild", canvas)

	w := &Window{
		Canvas: canvas,
		Width:  width,
		Height: height,
	}

	js.Global.Call("addEventListener", "keyup", func(event *js.Object) {
		e := &KeyEvent{
			KeyCode: event.Get("keyCode").Int(),
		}

		for _, l := range w.keyListeners {
			l(e)
		}

	}, false)

	canvas.Call("addEventListener", "click", func(event *js.Object) {
		e := &MouseEvent{
			Pos: num.Vec2{
				X: float32(event.Get("clientX").Float()),
				Y: float32(event.Get("clientY").Float()),
			},
		}

		for _, l := range w.mouseListeners {
			l(e)
		}

	}, false)

	return w
}

func (w *Window) AddKeyListener(l KeyListener) {
	w.keyListeners = append(w.keyListeners, l)
}

func (w *Window) AddMouseListener(l MouseListener) {
	w.mouseListeners = append(w.mouseListeners, l)
}
