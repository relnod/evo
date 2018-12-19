package math32

import "fmt"

type Mat4 struct {
	Data []float32
}

func NewMat4(x1, x2, x3, x4, y1, y2, y3, y4, z1, z2, z3, z4, w1, w2, w3, w4 float32) *Mat4 {
	m := Mat4{}
	m.Data = make([]float32, 16, 16)

	m.SetData(x1, x2, x3, x4, y1, y2, y3, y4, z1, z2, z3, z4, w1, w2, w3, w4)

	return &m
}

func NewMat4Empty() *Mat4 {
	m := Mat4{}
	m.Data = make([]float32, 16, 16)

	return &m
}

func (m *Mat4) SetData(x1, x2, x3, x4, y1, y2, y3, y4, z1, z2, z3, z4, w1, w2, w3, w4 float32) {
	m.Data[0] = x1
	m.Data[1] = x2
	m.Data[2] = x3
	m.Data[3] = x4

	m.Data[4] = y1
	m.Data[5] = y2
	m.Data[6] = y3
	m.Data[7] = y4

	m.Data[8] = z1
	m.Data[9] = z2
	m.Data[10] = z3
	m.Data[11] = z4

	m.Data[12] = w1
	m.Data[13] = w2
	m.Data[14] = w3
	m.Data[15] = w4
}

func (m1 *Mat4) Mult(m2 *Mat4) *Mat4 {
	// m := NewMat4Empty()
	// for r := 0; r < 4; r++ {
	// 	for c := 0; c < 4; c++ {
	// 		var val float32

	// 		for i := 0; i < 4; i++ {
	// 			val += m1.Data[r*4+i] * m2.Data[c+4*i]
	// 		}

	// 		m.Data[r*4+c] = val
	// 	}
	// }

	return NewMat4(
		m1.Data[0*4+0]*m2.Data[0+4*0]+
			m1.Data[0*4+1]*m2.Data[0+4*1]+
			m1.Data[0*4+2]*m2.Data[0+4*2]+
			m1.Data[0*4+3]*m2.Data[0+4*3],
		m1.Data[0*4+0]*m2.Data[1+4*0]+
			m1.Data[0*4+1]*m2.Data[1+4*1]+
			m1.Data[0*4+2]*m2.Data[1+4*2]+
			m1.Data[0*4+3]*m2.Data[1+4*3],
		m1.Data[0*4+0]*m2.Data[2+4*0]+
			m1.Data[0*4+1]*m2.Data[2+4*1]+
			m1.Data[0*4+2]*m2.Data[2+4*2]+
			m1.Data[0*4+3]*m2.Data[2+4*3],
		m1.Data[0*4+0]*m2.Data[3+4*0]+
			m1.Data[0*4+1]*m2.Data[3+4*1]+
			m1.Data[0*4+2]*m2.Data[3+4*2]+
			m1.Data[0*4+3]*m2.Data[3+4*3],

		m1.Data[1*4+0]*m2.Data[0+4*0]+
			m1.Data[1*4+1]*m2.Data[0+4*1]+
			m1.Data[1*4+2]*m2.Data[0+4*2]+
			m1.Data[1*4+3]*m2.Data[0+4*3],
		m1.Data[1*4+0]*m2.Data[1+4*0]+
			m1.Data[1*4+1]*m2.Data[1+4*1]+
			m1.Data[1*4+2]*m2.Data[1+4*2]+
			m1.Data[1*4+3]*m2.Data[1+4*3],
		m1.Data[1*4+0]*m2.Data[2+4*0]+
			m1.Data[1*4+1]*m2.Data[2+4*1]+
			m1.Data[1*4+2]*m2.Data[2+4*2]+
			m1.Data[1*4+3]*m2.Data[2+4*3],
		m1.Data[1*4+0]*m2.Data[3+4*0]+
			m1.Data[1*4+1]*m2.Data[3+4*1]+
			m1.Data[1*4+2]*m2.Data[3+4*2]+
			m1.Data[1*4+3]*m2.Data[3+4*3],

		m1.Data[2*4+0]*m2.Data[0+4*0]+
			m1.Data[2*4+1]*m2.Data[0+4*1]+
			m1.Data[2*4+2]*m2.Data[0+4*2]+
			m1.Data[2*4+3]*m2.Data[0+4*3],
		m1.Data[2*4+0]*m2.Data[1+4*0]+
			m1.Data[2*4+1]*m2.Data[1+4*1]+
			m1.Data[2*4+2]*m2.Data[1+4*2]+
			m1.Data[2*4+3]*m2.Data[1+4*3],
		m1.Data[2*4+0]*m2.Data[2+4*0]+
			m1.Data[2*4+1]*m2.Data[2+4*1]+
			m1.Data[2*4+2]*m2.Data[2+4*2]+
			m1.Data[2*4+3]*m2.Data[2+4*3],
		m1.Data[2*4+0]*m2.Data[3+4*0]+
			m1.Data[2*4+1]*m2.Data[3+4*1]+
			m1.Data[2*4+2]*m2.Data[3+4*2]+
			m1.Data[2*4+3]*m2.Data[3+4*3],

		m1.Data[3*4+0]*m2.Data[0+4*0]+
			m1.Data[3*4+1]*m2.Data[0+4*1]+
			m1.Data[3*4+2]*m2.Data[0+4*2]+
			m1.Data[3*4+3]*m2.Data[0+4*3],
		m1.Data[3*4+0]*m2.Data[1+4*0]+
			m1.Data[3*4+1]*m2.Data[1+4*1]+
			m1.Data[3*4+2]*m2.Data[1+4*2]+
			m1.Data[3*4+3]*m2.Data[1+4*3],
		m1.Data[3*4+0]*m2.Data[2+4*0]+
			m1.Data[3*4+1]*m2.Data[2+4*1]+
			m1.Data[3*4+2]*m2.Data[2+4*2]+
			m1.Data[3*4+3]*m2.Data[2+4*3],
		m1.Data[3*4+0]*m2.Data[3+4*0]+
			m1.Data[3*4+1]*m2.Data[3+4*1]+
			m1.Data[3*4+2]*m2.Data[3+4*2]+
			m1.Data[3*4+3]*m2.Data[3+4*3],
	)
}

func (m1 *Mat4) Transpose() *Mat4 {
	m := NewMat4Empty()
	for r := 0; r < 4; r++ {
		for c := 0; c < 4; c++ {
			m.Data[r*4+c] = m1.Data[r+4*c]
		}
	}

	return m
}

func (m *Mat4) String() string {
	str := ""
	for r := 0; r < 16; r += 4 {
		str += fmt.Sprintf("[ %f %f %f %f]\n", m.Data[r], m.Data[r+1], m.Data[r+2], m.Data[r+3])
	}

	return str
}
