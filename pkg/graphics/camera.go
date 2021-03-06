package graphics

import (
	"github.com/relnod/evo/pkg/math32"
)

// Renderer defines an interface for a render device.
type Renderer interface {
	SetViewport(x, y, width, height float32)
}

// Camera defines a 2D camera, that can zoom and move in all four directions.
type Camera struct {
	// window size
	windowWidth  int
	windowHeight int

	// world size
	worldWidth  int
	worldHeight int

	// zoom stores the zoom level.
	// Lowest value is 1.
	zoom int

	// x,y, width and height represent the viewport in world coordinates.
	x      int
	y      int
	width  int
	height int

	renderer Renderer
}

// NewCamera returns a new camera with a given window size.
func NewCamera(width, height int) *Camera {
	return &Camera{
		windowWidth:  width,
		windowHeight: height,

		worldWidth:  width,
		worldHeight: height,

		zoom: 1,

		x: 0,
		y: 0,

		width:  width,
		height: height,
	}
}

// Connect connects the camera to a given renderer.
func (c *Camera) Connect(renderer Renderer) {
	c.renderer = renderer
}

// SetSize sets the window size.
func (c *Camera) SetSize(widht, height int) {
	c.windowWidth = widht
	c.windowHeight = height
}

func (c *Camera) relativeToWindowSize(x, y float32) (float32, float32) {
	return x / float32(c.windowWidth), y / float32(c.windowHeight)
}

// santizeBounds ensures, that the viewport is within the world boundaries
func (c *Camera) santizeBounds() {
	// x and y
	if c.x < 0 {
		c.x = 0
	}
	if c.y < 0 {
		c.y = 0
	}
	if c.x+c.width > c.worldWidth {
		c.x = c.worldWidth - c.width
	}
	if c.y+c.height > c.worldHeight {
		c.y = c.worldHeight - c.height
	}

	// width and height
	if c.width < 100 {
		c.width = 100
	}
	if c.height < 100 {
		c.height = 100
	}
	if c.width > c.worldWidth {
		c.width = c.worldWidth
	}
	if c.height > c.worldHeight {
		c.height = c.worldHeight
	}
}

func (c *Camera) updateZoom() {
	if c.zoom <= 0 {
		c.zoom = 1
	}
	c.width = c.worldWidth / c.zoom
	c.height = c.worldHeight / c.zoom
}

// Viewport returns the viewport of the camera.
// All retruned values are between 0 and 1.
func (c *Camera) Viewport() (x float32, y float32, width float32, height float32) {
	if c.x > 0 {
		x = float32(c.x) / float32(c.windowWidth)
	}
	if c.y > 0 {
		y = float32(c.y) / float32(c.worldHeight)
	}

	ratio := math32.Min(
		float32(c.windowWidth)/float32(c.worldWidth),
		float32(c.windowHeight)/float32(c.worldHeight),
	)
	width, height = c.relativeToWindowSize(ratio, ratio)
	width /= float32(c.width) / float32(c.worldWidth)
	height /= float32(c.height) / float32(c.worldHeight)
	return x, y, width, height
}

// Update updates the viewport of the renderer.
func (c *Camera) Update() {
	if c.renderer == nil {
		return
	}

	c.updateZoom()
	c.santizeBounds()

	c.renderer.SetViewport(c.Viewport())
}

// ZoomIn zooms the viewport in.
func (c *Camera) ZoomIn() {
	c.zoom++
	c.Update()
}

// ZoomOut zooms the viewport out.
func (c *Camera) ZoomOut() {
	c.zoom--
	c.Update()
}

// MoveDown moves the viewport down.
func (c *Camera) MoveDown() {
	c.y += c.worldHeight / 10
	c.Update()
}

// MoveUp moves the viewport up.
func (c *Camera) MoveUp() {
	c.y -= c.worldHeight / 10
	c.Update()
}

// MoveRight moves the viewport to the right.
func (c *Camera) MoveRight() {
	c.x += c.worldWidth / 10
	c.Update()
}

// MoveLeft moves the viewport to the left.
func (c *Camera) MoveLeft() {
	c.x -= c.worldWidth / 10
	c.Update()
}
