package vecmath

import(
	"testing"
    // . "github.com/r7kamura/gospel"
)

func TestNewVec(t *testing.T) {
	v := NewVecZero()
	v.Length()
}

