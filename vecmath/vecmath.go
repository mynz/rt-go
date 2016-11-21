package vecmath

import (
	"math"
)

type Vec struct {
	V [3]float32
}

func NewVecZero() Vec {
	return Vec{[3]float32{0.0, 0.0, 0.0}}
}

func (v Vec) Length() float32 {
	return float32(math.Sqrt(float64(v.V[0]*v.V[0] + v.V[1]*v.V[1] + v.V[2]*v.V[2])))
}
