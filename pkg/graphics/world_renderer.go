package graphics

import (
	"encoding/binary"
	"log"
	"math"

	"github.com/gopherjs/gopherjs/js"
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
	VB          gl.Buffer
	ItemSize    int
	math32Items int
}

type WorldRenderer struct {
	width  float32
	height float32

	viewportWidth  float32
	viewportHeight float32

	canvas *js.Object

	gl      *webgl.Context
	program gl.Program
	aspect  float32

	aVertexPosition gl.Attrib
	uColor          gl.Uniform
	mModel          gl.Uniform
	mWorld          gl.Uniform

	circle RenderType
}

func NewWorldRenderer(width, height float32) *WorldRenderer {
	return &WorldRenderer{
		width:  width,
		height: height,

		viewportWidth:  width,
		viewportHeight: height,
	}
}

func (r *WorldRenderer) Update(w *world.World) {
	r.Clear()

	for _, c := range w.Creatures {
		if c.Speed == 0 {
			r.SetColor(0.0, 1.0-4.0/c.Radius/3.0, 0.0, 0.0)
		} else {
			r.SetColor(1/(c.Radius-4.0), 0.0, 0.0, 1.0)
		}
		r.DrawCircle(c.Pos.X, c.Pos.Y, c.Radius)
	}

	if gl.GetError() != gl.NO_ERROR {
		log.Println("OPENGL Error: ", gl.GetString(gl.GetError()))
	}
}

func (r *WorldRenderer) SetSize(width, height int) {
	gl.Viewport(0, 0, width, height)
	r.width = float32(width)
	r.height = float32(height)
}

func (r *WorldRenderer) Init() {
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

	r.program = program
	r.aspect = r.width / r.height
	r.aVertexPosition = gl.GetAttribLocation(program, "aVertexPosition")
	r.uColor = gl.GetUniformLocation(program, "uColor")
	r.mModel = gl.GetUniformLocation(program, "mModel")
	r.mWorld = gl.GetUniformLocation(program, "mWorld")

	// gl.Get
	// gl.Setprogram.Set("uColor", r.uColor)

	r.initCircleType()
}

func (r WorldRenderer) UpdateViewport(zoom, x, y float32) {
	dw := r.width / r.viewportWidth
	dh := r.height / r.viewportHeight
	d := dw
	if dw > dh {
		d = dh
	}
	mScale := math32.NewMat4(
		d*2.0/r.width*zoom, 0, 0, -x,
		0, -d*2.0/r.height*zoom, 0, y,
		0, 0, 1, 0,
		0, 0, 0, 1,
	)
	mTranslation := math32.NewMat4(
		1, 0, 0, -1,
		0, 1, 0, 1,
		0, 0, 1, 0,
		0, 0, 0, 1,
	)

	gl.UniformMatrix4fv(r.mWorld, mTranslation.Mult(mScale).Transpose().Data)
}

func (r *WorldRenderer) initCircleType() {
	vertices := make([]float32, 400)
	vertices[0] = 0
	vertices[1] = 0

	resolution := 2 * math.Pi / float64(len(vertices)/2-2)
	s := 0.0
	for i := 2; i < len(vertices); i += 2 {
		vertices[i] = float32(math.Cos(s))
		vertices[i+1] = float32(math.Sin(s))

		s += resolution
	}

	vbuffer := gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, vbuffer)
	verticesBytes := f32.Bytes(binary.LittleEndian, vertices...)
	gl.BufferData(gl.ARRAY_BUFFER, verticesBytes, gl.STATIC_DRAW)
	// gl.BindBuffer(gl.ARRAY_BUFFER, nil)

	itemSize := 2
	math32Items := len(vertices) / itemSize

	r.circle = RenderType{VB: vbuffer, ItemSize: itemSize, math32Items: math32Items}
}

func (r *WorldRenderer) Clear() {
	gl.Clear(gl.COLOR_BUFFER_BIT)
}

func (w *WorldRenderer) SetColor(r, g, b, a float32) {
	gl.Uniform4f(w.uColor, r, g, b, a)
}

func (r *WorldRenderer) DrawCircle(x, y, radius float32) {
	gl.BindBuffer(gl.ARRAY_BUFFER, r.circle.VB)

	gl.EnableVertexAttribArray(r.aVertexPosition)
	gl.VertexAttribPointer(r.aVertexPosition, r.circle.ItemSize, gl.FLOAT, false, 0, 0)

	mScale := math32.NewMat4(
		radius, 0, 0, 0,
		0, radius, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	)
	mTranslation := math32.NewMat4(
		1, 0, 0, x,
		0, 1, 0, y,
		0, 0, 1, 0,
		0, 0, 0, 1,
	)
	gl.UniformMatrix4fv(r.mModel, mTranslation.Mult(mScale).Transpose().Data)

	gl.DrawArrays(gl.TRIANGLE_FAN, 0, r.circle.math32Items)
}
