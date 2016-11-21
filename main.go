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

func Vmul(s float32, v mgl32.Vec3) mgl32.Vec3 { return v.Mul(s) }
func Vadd(a, b mgl32.Vec3) mgl32.Vec3 { return a.Add(b) }
func Vsub(a, b mgl32.Vec3) mgl32.Vec3 { return a.Sub(b) }
func Vdot(a, b mgl32.Vec3) float32 { return a.Dot(b) }

////

type Ray struct {
	A, B mgl32.Vec3
}

func (r Ray) Origin() mgl32.Vec3 { return r.A }
func (r Ray) Direction() mgl32.Vec3 { return r.B }
func (r Ray) PointAtParameter(t float32) mgl32.Vec3 { return Vadd(r.A, Vmul(t, r.B)) }

////

func HitSphere(center mgl32.Vec3, radius float32, r Ray) bool {
	oc := Vsub(r.Origin(), center)
	a := Vdot(r.Direction(), r.Direction())
	b := 2.0 * Vdot(oc, r.Direction())
	c := Vdot(oc, oc) - radius * radius
	discriminant := b * b - 4 * a * c
	return discriminant > 0
}

func CalcColor(r Ray) mgl32.Vec3 {
	if ( HitSphere(mgl32.Vec3{0, 0, -1}, 0.5, r) ) {
		return mgl32.Vec3{1, 0, 0}
	}

	unitDirection := r.Direction().Normalize()
	t := 0.5 * (unitDirection.Y() + 1.0)
	return Vadd(Vmul(1.0 - t, mgl32.Vec3{1, 1, 1}), Vmul(t, mgl32.Vec3{0.5, 0.7, 1.0}))
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
			dir := Vadd(lowerLeftCorner, Vadd(Vmul(u, horizontal), Vmul(v, vertical)))
			ray := Ray{ origin, dir }
			col := ConvToColor(CalcColor(ray))
			img.Set(i, ny - j, col) // reverse height
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


