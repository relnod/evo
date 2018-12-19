package graphics

// Camera defines a 2D camera, that can zoom and move in all four directions.
type Camera struct {
	// zoom should alwys be above 1.0
	zoom float32

	// x should always be between 0.0 and 0.9
	x float32
	// y should always be between 0.0 and 0.9
	y float32

	renderer Renderer
}

// NewCamera returns a new camera.
func NewCamera(renderer Renderer) *Camera {
	return &Camera{
		zoom: 1.0,

		x: 0.0,
		y: 0.0,

		renderer: renderer,
	}
}

// Update updates the viewport of the renderer.
func (c *Camera) Update() {
	c.renderer.UpdateViewport(c.zoom, c.x, c.y)
}

// ZoomIn zooms the viewport in.
func (c *Camera) ZoomIn() {
	c.zoom += 0.5
	c.Update()
}

// ZoomOut zooms the viewport out.
func (c *Camera) ZoomOut() {
	if c.zoom <= 1.0 {
		c.zoom = 1.0
		return
	}
	c.zoom -= 0.5
	c.Update()
}

// MoveDown moves the viewport down.
func (c *Camera) MoveDown() {
	if c.zoom == 1.0 {
		return
	}
	if c.y >= 0.9 {
		return
	}
	c.y += 0.1
	c.Update()
}

// MoveUp moves the viewport up.
func (c *Camera) MoveUp() {
	if c.zoom == 1.0 {
		return
	}
	if c.y <= 0.0 {
		return
	}
	c.y -= 0.1
	c.Update()
}

// MoveRight moves the viewport to the right.
func (c *Camera) MoveRight() {
	if c.zoom == 1.0 {
		return
	}
	if c.x >= 0.9 {
		return
	}
	c.x += 0.1
	c.Update()
}

// MoveLeft moves the viewport to the left.
func (c *Camera) MoveLeft() {
	if c.zoom == 1.0 {
		return
	}
	if c.x <= 0.0 {
		return
	}
	c.x -= 0.1
	c.Update()
}
