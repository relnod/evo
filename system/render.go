package system

import (
	"log"
	"math"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/webgl"
	"github.com/relnod/evo/num"
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
	VB       *js.Object
	ItemSize int
	NumItems int
}

type Render struct {
	system *System

	canvas *js.Object

	gl      *webgl.Context
	program *js.Object
	aspect  float32

	aVertexPosition int
	uColor          *js.Object
	mModel          *js.Object
	mWorld          *js.Object

	circle RenderType
}

func NewRender(s *System, canvas *js.Object) *Render {
	return &Render{system: s, canvas: canvas}
}

func (r *Render) Update() {
	r.Clear()

	for _, c := range r.system.creatures {
		r.SetColor(1.0/c.Radius, 0.2, 1.0/c.Radius, 1.0)
		r.DrawCircle(c.Pos.X, c.Pos.Y, c.Radius)
	}
}

func (r *Render) Init() {
	attrs := webgl.DefaultAttributes()
	attrs.Alpha = false

	gl, err := webgl.NewContext(r.canvas, attrs)
	if err != nil {
		js.Global.Call("alert", "Error: "+err.Error())
	}

	gl.Viewport(0, 0, int(r.system.Width), int(r.system.Height))
	gl.ClearColor(1.0, 1.0, 1.0, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	vs := gl.CreateShader(gl.VERTEX_SHADER)
	gl.ShaderSource(vs, vertexShader)
	gl.CompileShader(vs)
	if !gl.GetShaderParameterb(vs, gl.COMPILE_STATUS) {
		log.Fatal(gl.GetShaderInfoLog(vs))
	}

	fs := gl.CreateShader(gl.FRAGMENT_SHADER)
	gl.ShaderSource(fs, fragmentShader)
	gl.CompileShader(fs)
	if !gl.GetShaderParameterb(fs, gl.COMPILE_STATUS) {
		log.Fatal(gl.GetShaderInfoLog(fs))
	}

	program := gl.CreateProgram()
	gl.AttachShader(program, vs)
	gl.AttachShader(program, fs)
	gl.LinkProgram(program)

	gl.UseProgram(program)

	mScale := num.NewMat4(
		2.0/r.system.Width, 0, 0, 0,
		0, -2.0/r.system.Height, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	)
	mTranslation := num.NewMat4(
		1, 0, 0, -1,
		0, 1, 0, 1,
		0, 0, 1, 0,
		0, 0, 0, 1,
	)

	r.gl = gl
	r.program = program
	r.aspect = r.system.Width / r.system.Height
	r.aVertexPosition = gl.GetAttribLocation(program, "aVertexPosition")
	r.uColor = gl.GetUniformLocation(program, "uColor")
	r.mModel = gl.GetUniformLocation(program, "mModel")
	r.mWorld = gl.GetUniformLocation(program, "mWorld")

	program.Set("uColor", r.uColor)

	gl.UniformMatrix4fv(r.mWorld, false, mTranslation.Mult(mScale).Transpose().Data)

	r.initCircleType()
}

func (r *Render) initCircleType() {
	gl := r.gl

	vertices := make([]float32, 200)
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
	gl.BufferData(gl.ARRAY_BUFFER, vertices, gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, nil)

	itemSize := 2
	numItems := len(vertices) / itemSize

	r.circle = RenderType{VB: vbuffer, ItemSize: itemSize, NumItems: numItems}
}

func (r *Render) Clear() {
	gl := r.gl

	gl.Clear(gl.COLOR_BUFFER_BIT)
}

func (Render *Render) SetColor(r, g, b, a float32) {
	Render.gl.Uniform4f(Render.uColor, r, g, b, a)
}

func (r *Render) DrawCircle(x, y, radius float32) {
	gl := r.gl
	program := r.program

	gl.BindBuffer(gl.ARRAY_BUFFER, r.circle.VB)

	program.Set("aVertexPosition", r.aVertexPosition)
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
	gl.UniformMatrix4fv(r.mModel, false, mTranslation.Mult(mScale).Transpose().Data)

	gl.DrawArrays(gl.TRIANGLE_FAN, 0, r.circle.NumItems)
	gl.BindBuffer(gl.ARRAY_BUFFER, nil)
}
