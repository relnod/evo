package graphics

import (
	"encoding/binary"
	"log"
	"math"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/webgl"
	"github.com/goxjs/gl"
	"golang.org/x/mobile/exp/f32"

	"github.com/relnod/evo/pkg/num"
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
	NumItems int
}

type Render struct {
	width  float32
	height float32

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

func NewRender(width, height float32) *Render {
	return &Render{width: width, height: height}
}

func (r *Render) UpdateWorld(w *world.World) {
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

func (r *Render) Init() {
	gl.Viewport(0, 0, int(r.width), int(r.height))
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

	mScale := num.NewMat4(
		2.0/r.width, 0, 0, 0,
		0, -2.0/r.height, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	)
	mTranslation := num.NewMat4(
		1, 0, 0, -1,
		0, 1, 0, 1,
		0, 0, 1, 0,
		0, 0, 0, 1,
	)

	r.program = program
	r.aspect = r.width / r.height
	r.aVertexPosition = gl.GetAttribLocation(program, "aVertexPosition")
	r.uColor = gl.GetUniformLocation(program, "uColor")
	r.mModel = gl.GetUniformLocation(program, "mModel")
	r.mWorld = gl.GetUniformLocation(program, "mWorld")

	// gl.Get
	// gl.Setprogram.Set("uColor", r.uColor)

	gl.UniformMatrix4fv(r.mWorld, mTranslation.Mult(mScale).Transpose().Data)

	r.initCircleType()
}

func (r *Render) initCircleType() {
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
	numItems := len(vertices) / itemSize

	r.circle = RenderType{VB: vbuffer, ItemSize: itemSize, NumItems: numItems}
}

func (r *Render) Clear() {
	gl.Clear(gl.COLOR_BUFFER_BIT)
}

func (Render *Render) SetColor(r, g, b, a float32) {
	gl.Uniform4f(Render.uColor, r, g, b, a)
}

func (r *Render) DrawCircle(x, y, radius float32) {
	// log.Println(x, y, radius)
	// program := r.program

	gl.BindBuffer(gl.ARRAY_BUFFER, r.circle.VB)

	gl.EnableVertexAttribArray(r.aVertexPosition)
	gl.VertexAttribPointer(r.aVertexPosition, r.circle.ItemSize, gl.FLOAT, false, 0, 0)

	mScale := num.NewMat4(
		radius, 0, 0, 0,
		0, radius, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	)
	mTranslation := num.NewMat4(
		1, 0, 0, x,
		0, 1, 0, y,
		0, 0, 1, 0,
		0, 0, 0, 1,
	)
	gl.UniformMatrix4fv(r.mModel, mTranslation.Mult(mScale).Transpose().Data)

	gl.DrawArrays(gl.TRIANGLE_FAN, 0, r.circle.NumItems)
	// gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}
