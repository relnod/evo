package math32

import "fmt"

type Mat4 [16]float32

func NewMat4(x1, x2, x3, x4, y1, y2, y3, y4, z1, z2, z3, z4, w1, w2, w3, w4 float32) Mat4 {
	return Mat4{x1, x2, x3, x4, y1, y2, y3, y4, z1, z2, z3, z4, w1, w2, w3, w4}
}

func NewMat4Empty() Mat4 {
	return Mat4{}
}

func (m Mat4) Data() []float32 {
	return m[:]
}

func (m Mat4) SetData(x1, x2, x3, x4, y1, y2, y3, y4, z1, z2, z3, z4, w1, w2, w3, w4 float32) {
	m[0] = x1
	m[1] = x2
	m[2] = x3
	m[3] = x4

	m[4] = y1
	m[5] = y2
	m[6] = y3
	m[7] = y4

	m[8] = z1
	m[9] = z2
	m[10] = z3
	m[11] = z4

	m[12] = w1
	m[13] = w2
	m[14] = w3
	m[15] = w4
}

func (m Mat4) Mult(m2 Mat4) Mat4 {
	// m := NewMat4Empty()
	// for r := 0; r < 4; r++ {
	// 	for c := 0; c < 4; c++ {
	// 		var val float32

	// 		for i := 0; i < 4; i++ {
	// 			val += m[r*4+i] * m2[c+4*i]
	// 		}

	// 		m[r*4+c] = val
	// 	}
	// }

	return NewMat4(
		m[0*4+0]*m2[0+4*0]+
			m[0*4+1]*m2[0+4*1]+
			m[0*4+2]*m2[0+4*2]+
			m[0*4+3]*m2[0+4*3],
		m[0*4+0]*m2[1+4*0]+
			m[0*4+1]*m2[1+4*1]+
			m[0*4+2]*m2[1+4*2]+
			m[0*4+3]*m2[1+4*3],
		m[0*4+0]*m2[2+4*0]+
			m[0*4+1]*m2[2+4*1]+
			m[0*4+2]*m2[2+4*2]+
			m[0*4+3]*m2[2+4*3],
		m[0*4+0]*m2[3+4*0]+
			m[0*4+1]*m2[3+4*1]+
			m[0*4+2]*m2[3+4*2]+
			m[0*4+3]*m2[3+4*3],

		m[1*4+0]*m2[0+4*0]+
			m[1*4+1]*m2[0+4*1]+
			m[1*4+2]*m2[0+4*2]+
			m[1*4+3]*m2[0+4*3],
		m[1*4+0]*m2[1+4*0]+
			m[1*4+1]*m2[1+4*1]+
			m[1*4+2]*m2[1+4*2]+
			m[1*4+3]*m2[1+4*3],
		m[1*4+0]*m2[2+4*0]+
			m[1*4+1]*m2[2+4*1]+
			m[1*4+2]*m2[2+4*2]+
			m[1*4+3]*m2[2+4*3],
		m[1*4+0]*m2[3+4*0]+
			m[1*4+1]*m2[3+4*1]+
			m[1*4+2]*m2[3+4*2]+
			m[1*4+3]*m2[3+4*3],

		m[2*4+0]*m2[0+4*0]+
			m[2*4+1]*m2[0+4*1]+
			m[2*4+2]*m2[0+4*2]+
			m[2*4+3]*m2[0+4*3],
		m[2*4+0]*m2[1+4*0]+
			m[2*4+1]*m2[1+4*1]+
			m[2*4+2]*m2[1+4*2]+
			m[2*4+3]*m2[1+4*3],
		m[2*4+0]*m2[2+4*0]+
			m[2*4+1]*m2[2+4*1]+
			m[2*4+2]*m2[2+4*2]+
			m[2*4+3]*m2[2+4*3],
		m[2*4+0]*m2[3+4*0]+
			m[2*4+1]*m2[3+4*1]+
			m[2*4+2]*m2[3+4*2]+
			m[2*4+3]*m2[3+4*3],

		m[3*4+0]*m2[0+4*0]+
			m[3*4+1]*m2[0+4*1]+
			m[3*4+2]*m2[0+4*2]+
			m[3*4+3]*m2[0+4*3],
		m[3*4+0]*m2[1+4*0]+
			m[3*4+1]*m2[1+4*1]+
			m[3*4+2]*m2[1+4*2]+
			m[3*4+3]*m2[1+4*3],
		m[3*4+0]*m2[2+4*0]+
			m[3*4+1]*m2[2+4*1]+
			m[3*4+2]*m2[2+4*2]+
			m[3*4+3]*m2[2+4*3],
		m[3*4+0]*m2[3+4*0]+
			m[3*4+1]*m2[3+4*1]+
			m[3*4+2]*m2[3+4*2]+
			m[3*4+3]*m2[3+4*3],
	)
}

func (m Mat4) Transpose() Mat4 {
	m2 := NewMat4Empty()
	for r := 0; r < 4; r++ {
		for c := 0; c < 4; c++ {
			m2[r*4+c] = m[r+4*c]
		}
	}

	return m2
}

func (m Mat4) String() string {
	str := ""
	for r := 0; r < 16; r += 4 {
		str += fmt.Sprintf("[ %f %f %f %f]\n", m[r], m[r+1], m[r+2], m[r+3])
	}

	return str
}
