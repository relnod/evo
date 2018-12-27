package graphics

import (
	"encoding/binary"
	"log"
	"math"

	"github.com/gopherjs/webgl"
	"github.com/goxjs/gl"
	"golang.org/x/mobile/exp/f32"

	"github.com/relnod/evo/pkg/math32"
	"github.com/relnod/evo/pkg/world"
)

var vertexShader = `
attribute vec2 aVertexPosition;
uniform mat4 mModel;
uniform mat4 mWorld;
void main() {
	gl_Position = mWorld * mModel * vec4(aVertexPosition, 0.0, 1.0);
}
`

var fragmentShader = `
#ifdef GL_ES
precision highp float;
#endif
uniform vec4 uColor;
void main() {
	gl_FragColor = uColor;
}
`

type RenderType struct {
	VB       gl.Buffer
	ItemSize int
	numItems int
}

type WorldRenderer struct {
	width  int
	height int

	viewportWidth  int
	viewportHeight int

	gl      *webgl.Context
	program gl.Program

	aVertexPosition gl.Attrib
	uColor          gl.Uniform
	mModel          gl.Uniform
	mWorld          gl.Uniform

	circle RenderType
}

func NewWorldRenderer(width, height int) *WorldRenderer {
	return &WorldRenderer{
		width:  width,
		height: height,

		viewportWidth:  width,
		viewportHeight: height,
	}
}

func (w *WorldRenderer) Update(world *world.World) {
	w.Clear()

	for _, c := range world.Creatures {
		if c.Speed == 0 {
			w.SetColor(0.0, 1.0-4.0/c.Radius/3.0, 0.0, 0.0)
		} else {
			w.SetColor(1/(c.Radius-4.0), 0.0, 0.0, 1.0)
		}
		w.DrawCircle(c.Pos.X, c.Pos.Y, c.Radius)

		if c.Eye != nil {
			w.SetColor(0.0, 0.0, 0.0, 0.0)
			w.DrawPartialCircle(c.Pos.X, c.Pos.Y, c.Eye.Range, c.Eye.FOV, math.Atan2(c.Eye.Dir.Y, c.Eye.Dir.X))
		}
	}

	if gl.GetError() != gl.NO_ERROR {
		log.Println("OPENGL Error: ", gl.GetString(gl.GetError()))
	}
}

func (w *WorldRenderer) SetSize(width, height int) {
	gl.Viewport(0, 0, width, height)
	w.width = width
	w.height = height
}

func (w *WorldRenderer) Init() {
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	vs := gl.CreateShader(gl.VERTEX_SHADER)
	gl.ShaderSource(vs, vertexShader)
	gl.CompileShader(vs)
	if gl.GetShaderi(vs, gl.COMPILE_STATUS) == 0 {
		log.Fatal(gl.GetShaderInfoLog(vs))
	}

	fs := gl.CreateShader(gl.FRAGMENT_SHADER)
	gl.ShaderSource(fs, fragmentShader)
	gl.CompileShader(fs)
	if gl.GetShaderi(fs, gl.COMPILE_STATUS) == 0 {
		log.Fatal(gl.GetShaderInfoLog(fs))
	}

	program := gl.CreateProgram()
	gl.AttachShader(program, vs)
	gl.AttachShader(program, fs)
	gl.LinkProgram(program)

	gl.UseProgram(program)

	w.program = program
	w.aVertexPosition = gl.GetAttribLocation(program, "aVertexPosition")
	w.uColor = gl.GetUniformLocation(program, "uColor")
	w.mModel = gl.GetUniformLocation(program, "mModel")
	w.mWorld = gl.GetUniformLocation(program, "mWorld")

	// gl.Get
	// gl.Setprogram.Set("uColor", r.uColor)

	w.initCircleType()
}

func (w WorldRenderer) UpdateViewport(zoom, x, y float64) {
	dw := float64(w.width) / float64(w.viewportWidth)
	dh := float64(w.height) / float64(w.viewportHeight)
	d := dw
	if dw > dh {
		d = dh
	}
	mScale := math32.NewMat4(
		float32(d*2.0/float64(w.width)*zoom), 0, 0, float32(-x),
		0, float32(-d*2.0/float64(w.height)*zoom), 0, float32(y),
		0, 0, 1, 0,
		0, 0, 0, 1,
	)
	mTranslation := math32.NewMat4(
		1, 0, 0, -1,
		0, 1, 0, 1,
		0, 0, 1, 0,
		0, 0, 0, 1,
	)

	gl.UniformMatrix4fv(w.mWorld, mTranslation.Mult(mScale).Transpose().Data())
}

func (w *WorldRenderer) initCircleType() {
	numVertices := 40
	vertices := make([]float32, numVertices*2)

	for i := 0; i < len(vertices); i += 2 {
		theta := 2.0 * math.Pi * (float64(i) / 2.0) / float64(numVertices)
		vertices[i] = float32(math.Cos(theta))
		vertices[i+1] = float32(math.Sin(theta))
	}

	vbuffer := gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, vbuffer)
	verticesBytes := f32.Bytes(binary.LittleEndian, vertices...)
	gl.BufferData(gl.ARRAY_BUFFER, verticesBytes, gl.STATIC_DRAW)

	itemSize := 2
	math32Items := len(vertices) / itemSize

	w.circle = RenderType{VB: vbuffer, ItemSize: itemSize, numItems: math32Items}
}

func (w *WorldRenderer) Clear() {
	gl.Clear(gl.COLOR_BUFFER_BIT)
}

func (w *WorldRenderer) SetColor(r, g, b, a float64) {
	gl.Uniform4f(w.uColor, float32(r), float32(g), float32(b), float32(a))
}

func (w *WorldRenderer) DrawCircle(x, y, radius float64) {
	gl.BindBuffer(gl.ARRAY_BUFFER, w.circle.VB)

	gl.EnableVertexAttribArray(w.aVertexPosition)
	gl.VertexAttribPointer(w.aVertexPosition, w.circle.ItemSize, gl.FLOAT, false, 0, 0)

	mScale := math32.NewMat4(
		float32(radius), 0, 0, 0,
		0, float32(radius), 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	)
	mTranslation := math32.NewMat4(
		1, 0, 0, float32(x),
		0, 1, 0, float32(y),
		0, 0, 1, 0,
		0, 0, 0, 1,
	)
	gl.UniformMatrix4fv(w.mModel, mTranslation.Mult(mScale).Transpose().Data())

	gl.DrawArrays(gl.TRIANGLE_FAN, 0, w.circle.numItems)
}

func (w *WorldRenderer) DrawPartialCircle(x, y, radius, fov, angle float64) {
	gl.BindBuffer(gl.ARRAY_BUFFER, w.circle.VB)

	gl.EnableVertexAttribArray(w.aVertexPosition)
	gl.VertexAttribPointer(w.aVertexPosition, w.circle.ItemSize, gl.FLOAT, false, 0, 0)

	mScale := math32.NewMat4(
		float32(radius), 0, 0, 0,
		0, float32(radius), 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	)
	mTranslation := math32.NewMat4(
		1, 0, 0, float32(x),
		0, 1, 0, float32(y),
		0, 0, 1, 0,
		0, 0, 0, 1,
	)
	// mRotation := math32.NewMat4(
	// 	1, 0, 0, 0,
	// 	0, float32(math.Cos(angle)), -1*float32(math.Sin(angle)), 0,
	// 	0, float32(math.Sin(angle)), float32(math.Cos(angle)), 0,
	// 	0, 0, 0, 1,
	// )
	// mRotation := math32.NewMat4(
	// 	float32(math.Cos(angle)), 0, float32(math.Sin(angle)), 0,
	// 	0, 1, 0, 0,
	// 	-1*float32(math.Sin(angle)), 0, float32(math.Cos(angle)), 0,
	// 	0, 0, 0, 1,
	// )
	mRotation := math32.NewMat4(
		float32(math.Cos(angle)), -1*float32(math.Sin(angle)), 0, 0,
		float32(math.Sin(angle)), float32(math.Cos(angle)), 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	)
	gl.UniformMatrix4fv(w.mModel, mTranslation.Mult(mScale).Mult(mRotation).Transpose().Data())

	gl.DrawArrays(gl.LINE_STRIP, 0, int(float64(w.circle.numItems)*fov/(2.0*math.Pi)))
}
