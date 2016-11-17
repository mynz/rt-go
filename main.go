package main

import(
	"os"
	"fmt"
	"image"
	"image/color"
	"image/png"

	// "./vecmath"
	"github.com/go-gl/mathgl/mgl32"
)

////

type Ray struct {
	A, B mgl32.Vec3
}

func (r Ray) Origin() mgl32.Vec3 { return r.A }
func (r Ray) Direction() mgl32.Vec3 { return r.B }
func (r Ray) PointAtParameter(t float32) mgl32.Vec3 { return r.A.Add(r.B.Mul(t)) }

////

func CalcColor(r Ray) mgl32.Vec3 {
	unitDirection := r.Direction().Normalize()
	t := 0.5 * (unitDirection.Y() + 1.0)
	return mgl32.Vec3{1, 1, 1}.Mul(1.0 - t).Add(mgl32.Vec3{0.5, 0.7, 1.0}.Mul(t))
}

func ConvToColor(c32 mgl32.Vec3) color.NRGBA {
	v := c32.Mul(255.99)
	return color.NRGBA{ uint8(v.X()), uint8(v.Y()), uint8(v.Z()), 255}
}


func RenderImage() image.Image {
	nx, ny := 200, 100

	lowerLeftCorner := mgl32.Vec3{ -2.0, -1.0, -1.0 }
	horizontal      := mgl32.Vec3{ 4.0, 0, 0 }
	vertical        := mgl32.Vec3{ 0, 2.0, 0 }
	origin          := mgl32.Vec3{ 0, 0, 0 }

	img := image.NewRGBA(image.Rect(0, 0, nx, ny))

	for j := ny-1; j >= 0; j-- {
		for  i := 0; i < nx; i++ {
			u, v := float32(i) / float32(nx), float32(j) / float32(ny)
			dir := lowerLeftCorner.Add(horizontal.Mul(u).Add(vertical.Mul(v)))
			ray := Ray{ origin, dir }
			col := ConvToColor(CalcColor(ray))
			img.Set(i, j, col)
		}
	}

	return img
}

func writePngFile(filename string, img image.Image) {
	fp, err := os.Create(filename)
	defer fp.Close()
	if err != nil {
		panic(err)
	}

	if err := png.Encode(fp, img); err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println("Hello.")

	{
		img := RenderImage()
		writePngFile("test.png", img)
	}

	fmt.Println("Done.")
}


