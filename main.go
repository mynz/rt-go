package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"

	"github.com/go-gl/mathgl/mgl32"
)

////

func Sqrt32(f float32) float32 { return float32(math.Sqrt(float64(f))) }

func Vmul(s float32, v mgl32.Vec3) mgl32.Vec3 { return v.Mul(s) }
func VDiv(v mgl32.Vec3, s float32) mgl32.Vec3 { return mgl32.Vec3{v.X() / s, v.Y() / s, v.Z() / s} }
func Vadd(a, b mgl32.Vec3) mgl32.Vec3         { return a.Add(b) }
func Vsub(a, b mgl32.Vec3) mgl32.Vec3         { return a.Sub(b) }
func Vdot(a, b mgl32.Vec3) float32            { return a.Dot(b) }

////

type Ray struct {
	A, B mgl32.Vec3
}

func (r Ray) Origin() mgl32.Vec3                    { return r.A }
func (r Ray) Direction() mgl32.Vec3                 { return r.B }
func (r Ray) PointAtParameter(t float32) mgl32.Vec3 { return Vadd(r.A, Vmul(t, r.B)) }

////

type HitRecord struct {
	T      float32
	P      mgl32.Vec3
	Normal mgl32.Vec3
}

type Hitable interface {
	Hit(ray Ray, tmin, tmax float32, rec *HitRecord) bool
}

////

type Sphere struct {
	Center mgl32.Vec3
	Radius float32
}

func (s Sphere) Hit(ray Ray, tmin, tmax float32, rec *HitRecord) bool {
	center := s.Center
	radius := s.Radius

	oc := Vsub(ray.Origin(), center)
	a := Vdot(ray.Direction(), ray.Direction())
	b := Vdot(oc, ray.Direction())
	c := Vdot(oc, oc) - radius*radius
	discriminant := b*b - a*c
	if discriminant > 0 {
		var tmp float32
		tmp = (-b - Sqrt32(b*b-a*c)) / a
		if tmp < tmax && tmp > tmin {
			rec.T = tmp
			rec.P = ray.PointAtParameter(rec.T)
			rec.Normal = VDiv(Vsub(rec.P, center), radius)
			return true
		}
		tmp = (-b + Sqrt32(b*b-a*c)) / a
		if tmp < tmax && tmp > tmin {
			rec.T = tmp
			rec.P = ray.PointAtParameter(rec.T)
			rec.Normal = VDiv(Vsub(rec.P, center), radius)
			return true
		}
	}

	return false
}

////

type HitableList struct {
	List []Hitable
}

func (hlist HitableList) Hit(ray Ray, tmin, tmax float32, rec *HitRecord) bool {
	tmpRec := HitRecord{}
	hitAnything := false
	closestSoFar := tmax
	for _, v := range hlist.List {
		if v.Hit(ray, tmin, closestSoFar, &tmpRec) {
			hitAnything = true
			closestSoFar = tmpRec.T
			*rec = tmpRec
		}
	}
	return hitAnything
}

////

type Camera struct {
	origin          mgl32.Vec3
	lowerLeftCorner mgl32.Vec3
	horizontal      mgl32.Vec3
	vertical        mgl32.Vec3
}

func NewCamera() Camera {
	return Camera{
		lowerLeftCorner: mgl32.Vec3{-2.0, -1.0, -1.0},
		horizontal:      mgl32.Vec3{4.0, 0, 0},
		vertical:        mgl32.Vec3{0, 2.0, 0},
		origin:          mgl32.Vec3{0, 0, 0},
	}
}

func (cam Camera) GetRay(u, v float32) Ray {
	dir := Vadd(cam.lowerLeftCorner, Vadd(Vmul(u, cam.horizontal), Vmul(v, cam.vertical)))
	return Ray{cam.origin, dir}
}

////

// color
func CalcColor(r Ray, world Hitable) mgl32.Vec3 {
	rec := HitRecord{}
	var MAXFLOAT = float32(100000000.0)
	if world.Hit(r, 0.0, MAXFLOAT, &rec) {
		return Vmul(0.5, Vadd(rec.Normal, mgl32.Vec3{1, 1, 1}))
	} else {
		unitDirection := r.Direction().Normalize()
		t := 0.5 * (unitDirection.Y() + 1.0)
		return Vadd(Vmul(1.0-t, mgl32.Vec3{1, 1, 1}), Vmul(t, mgl32.Vec3{0.5, 0.7, 1.0}))
	}
}

func ConvToColor(c32 mgl32.Vec3) color.NRGBA {
	x, y, z := c32.Mul(255.99).Elem()
	return color.NRGBA{uint8(x), uint8(y), uint8(z), 255}
}

// main function.
func RenderImage() image.Image {
	nx, ny := 200, 100

	lowerLeftCorner := mgl32.Vec3{-2.0, -1.0, -1.0}
	horizontal := mgl32.Vec3{4.0, 0, 0}
	vertical := mgl32.Vec3{0, 2.0, 0}
	origin := mgl32.Vec3{0, 0, 0}

	world := HitableList{
		List: []Hitable{
			Sphere{mgl32.Vec3{0, 0, -1}, 0.5},
			Sphere{mgl32.Vec3{0, -100.5, -1}, 100.0},
		},
	}

	img := image.NewRGBA(image.Rect(0, 0, nx, ny))

	for j := ny - 1; j >= 0; j-- {
		for i := 0; i < nx; i++ {
			u, v := float32(i)/float32(nx), float32(j)/float32(ny)
			dir := Vadd(lowerLeftCorner, Vadd(Vmul(u, horizontal), Vmul(v, vertical)))
			ray := Ray{origin, dir}
			col := ConvToColor(CalcColor(ray, world))
			img.Set(i, ny-j, col) // reverse height
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
