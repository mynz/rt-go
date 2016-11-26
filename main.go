package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"

	"github.com/go-gl/mathgl/mgl32"
)

////

func Sqrt32(f float32) float32  { return float32(math.Sqrt(float64(f))) }
func Fpow(x, y float32) float32 { return float32(math.Pow(float64(x), float64(y))) }

func Vmul(s float32, v mgl32.Vec3) mgl32.Vec3 { return v.Mul(s) }
func Vdiv(v mgl32.Vec3, s float32) mgl32.Vec3 { return mgl32.Vec3{v.X() / s, v.Y() / s, v.Z() / s} }
func Vadd(a, b mgl32.Vec3) mgl32.Vec3         { return a.Add(b) }
func Vsub(a, b mgl32.Vec3) mgl32.Vec3         { return a.Sub(b) }
func Vdot(a, b mgl32.Vec3) float32            { return a.Dot(b) }

func VmulPerElem(u mgl32.Vec3, v mgl32.Vec3) mgl32.Vec3 {
	return mgl32.Vec3{u.X() * v.X(), u.Y() * v.Y(), u.Z() * v.Z()}
}

////

type Ray struct {
	A, B mgl32.Vec3
}

func (r Ray) Origin() mgl32.Vec3                    { return r.A }
func (r Ray) Direction() mgl32.Vec3                 { return r.B }
func (r Ray) PointAtParameter(t float32) mgl32.Vec3 { return Vadd(r.A, Vmul(t, r.B)) }

////

type Material interface {
	Scatter(ray Ray, rec HitRecord) (b bool, attenuation mgl32.Vec3, scattered Ray)
}

type Lambertian struct {
	albedo mgl32.Vec3
}

func NewLambertian(a mgl32.Vec3) *Lambertian { return &Lambertian{albedo: a} }

func (lam Lambertian) Scatter(ray Ray, rec HitRecord) (bool, mgl32.Vec3, Ray) {
	target := Vadd(Vadd(rec.P, rec.Normal), RandomInUnitSphere())
	scattered := Ray{rec.P, Vsub(target, rec.P)}
	attenuation := lam.albedo
	return true, attenuation, scattered
}

type Metal struct {
	albedo mgl32.Vec3
	fuzz   float32
}

func NewMetal(a mgl32.Vec3, f float32) *Metal { return &Metal{a, f} }

func reflect(v, n mgl32.Vec3) mgl32.Vec3 {
	return Vsub(v, Vmul(2, Vmul(Vdot(v, n), n)))
}

func (met Metal) Scatter(ray Ray, rec HitRecord) (bool, mgl32.Vec3, Ray) {
	reflected := reflect(ray.Direction().Normalize(), rec.Normal)
	// scattered := Ray{rec.P, reflected}
	scattered := Ray{rec.P, Vadd(reflected, Vmul(met.fuzz, RandomInUnitSphere()))}
	attenuation := met.albedo
	b := Vdot(scattered.Direction(), rec.Normal) > 0.0
	return b, attenuation, scattered
}

type Dielectric struct {
	refIdx float32
}

func NewDielectric(refIdx float32) *Dielectric { return &Dielectric{refIdx} }

func refract(v, n mgl32.Vec3, niOverNt float32) (bool, mgl32.Vec3) {
	uv := v.Normalize()
	dt := Vdot(uv, n)
	discriminant := 1.0 - niOverNt*niOverNt*(1-dt*dt)
	if discriminant > 0 {
		refracted := Vsub(Vmul(niOverNt, Vsub(uv, Vmul(dt, n))), Vmul(Sqrt32(discriminant), n))
		return true, refracted
	}
	return false, mgl32.Vec3{0, 0, 0}
}

func schlick(cosine, refIdx float32) float32 {
	r0 := (1 - refIdx) / (1 + refIdx)
	r0 = r0 * r0
	return r0 + (1-r0)*Fpow((1-cosine), 5)
}

func (die Dielectric) Scatter(ray Ray, rec HitRecord) (bool, mgl32.Vec3, Ray) {
	var outwardNormal mgl32.Vec3
	reflected := reflect(ray.Direction(), rec.Normal)
	var niOverNt float32
	attenuation := mgl32.Vec3{1.0, 1.0, 1.0}

	var reflectedProb, cosine float32

	if Vdot(ray.Direction(), rec.Normal) > 0 {
		outwardNormal = Vmul(-1.0, rec.Normal)
		niOverNt = die.refIdx
		cosine = die.refIdx * Vdot(ray.Direction(), rec.Normal) / ray.Direction().Len()
	} else {
		outwardNormal = rec.Normal
		niOverNt = 1.0 / die.refIdx
		cosine = -Vdot(ray.Direction(), rec.Normal) / ray.Direction().Len()
	}

	b, refracted := refract(ray.Direction(), outwardNormal, niOverNt)
	if b {
		reflectedProb = schlick(cosine, die.refIdx)
	} else {
		reflectedProb = 1.0
	}

	var scattered Ray
	if rand.Float32() < reflectedProb {
		scattered = Ray{rec.P, reflected}
	} else {
		scattered = Ray{rec.P, refracted}
	}
	return true, attenuation, scattered
}

////

type HitRecord struct {
	T      float32
	P      mgl32.Vec3
	Normal mgl32.Vec3
	MatPtr Material
}

type Hitable interface {
	Hit(ray Ray, tmin, tmax float32, rec *HitRecord) bool
}

////

type Sphere struct {
	Center mgl32.Vec3
	Radius float32
	MatPtr Material
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
			rec.Normal = Vdiv(Vsub(rec.P, center), radius)
			rec.MatPtr = s.MatPtr
			return true
		}
		tmp = (-b + Sqrt32(b*b-a*c)) / a
		if tmp < tmax && tmp > tmin {
			rec.T = tmp
			rec.P = ray.PointAtParameter(rec.T)
			rec.Normal = Vdiv(Vsub(rec.P, center), radius)
			rec.MatPtr = s.MatPtr
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
	u, v, w         mgl32.Vec3
	lensRadius      float32
}

func NewCamera(lookFrom, lookAt, vup mgl32.Vec3, vfov, aspect, aperture, focusDist float32) Camera {
	lensRadius := aperture / 2.0
	theta := float64(vfov * math.Pi / 180.0)
	halfHeight := float32(math.Tan(theta / 2.0))
	halfWidth := float32(aspect * halfHeight)
	origin := lookFrom
	w := Vsub(lookFrom, lookAt).Normalize()
	u := vup.Cross(w).Normalize()
	v := w.Cross(u)

	a0 := origin
	a1 := Vmul(halfWidth*focusDist, u)
	a2 := Vmul(halfHeight*focusDist, v)
	a3 := Vmul(focusDist, w)
	lowerLeftCorner := Vsub(Vsub(Vsub(a0, a1), a2), a3)

	horizontal := Vmul(2*halfWidth*focusDist, u)
	vertical := Vmul(2*halfHeight*focusDist, v)
	return Camera{
		lowerLeftCorner: lowerLeftCorner,
		horizontal:      horizontal,
		vertical:        vertical,
		origin:          origin,
		u:               u,
		v:               v,
		w:               w,
		lensRadius:      lensRadius,
	}
}

func randomInUnitDisk() mgl32.Vec3 {
	var p mgl32.Vec3
	for {
		p = Vsub(Vmul(2, mgl32.Vec3{rand.Float32(), rand.Float32(), 0}), mgl32.Vec3{1, 1, 0})
		if !(Vdot(p, p) >= 1.0) {
			break
		}
	}
	return p
}

func (cam Camera) GetRay(s, t float32) Ray {
	rd := Vmul(cam.lensRadius, randomInUnitDisk())
	offset := Vadd(Vmul(rd.X(), cam.u), Vmul(rd.Y(), cam.v))

	d0 := cam.lowerLeftCorner
	d1 := Vmul(s, cam.horizontal)
	d2 := Vmul(t, cam.vertical)
	d3 := cam.origin
	d4 := offset
	dir := Vadd(d0, Vadd(d1, d2))
	dir = Vsub(Vsub(dir, d3), d4)
	return Ray{Vadd(cam.origin, offset), dir}
}

////

func RandomInUnitSphere() mgl32.Vec3 {
	var p mgl32.Vec3
	for {
		p = Vsub(Vmul(2.0, mgl32.Vec3{rand.Float32(), rand.Float32(), rand.Float32()}), mgl32.Vec3{1, 1, 1})
		// if Vdot(p, p) < 1.0 {
		if !(Vdot(p, p) >= 1.0) {
			break
		}
	}
	return p
}

// color
func CalcColor(r Ray, world Hitable, depth int) mgl32.Vec3 {
	rec := HitRecord{}
	var MAXFLOAT = float32(100000000.0)
	if world.Hit(r, 0.001, MAXFLOAT, &rec) {
		if depth < 50 {
			if b, attenuation, scattered := rec.MatPtr.Scatter(r, rec); b {
				return VmulPerElem(attenuation, CalcColor(scattered, world, depth+1))
			}
		}
		return mgl32.Vec3{0, 0, 0}
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
	// nx, ny := 640, 480
	// ns := 100
	ns := 10 // original: 100

	// R := float32(math.Cos(math.Pi / 4))

	lookFrom := mgl32.Vec3{3, 3, 2}
	lookAt := mgl32.Vec3{0, 0, -1}
	distToFocus := Vsub(lookFrom, lookAt).Len()
	aperture := float32(2.0)

	cam := NewCamera(lookFrom, lookAt, mgl32.Vec3{0, 1, 0}, 20.0, float32(nx)/float32(ny), aperture, distToFocus)
	world := HitableList{
		List: []Hitable{

			/*
			 * Sphere{mgl32.Vec3{-R, 0, -1}, R, NewLambertian(mgl32.Vec3{0, 0, 1})},
			 * Sphere{mgl32.Vec3{+R, 0, -1}, R, NewLambertian(mgl32.Vec3{1, 0, 0})},
			 */

			Sphere{mgl32.Vec3{0, 0, -1}, 0.5, NewLambertian(mgl32.Vec3{0.1, 0.2, 0.5})},
			Sphere{mgl32.Vec3{0, -100.5, -1}, 100.0, NewLambertian(mgl32.Vec3{0.8, 0.8, 0.0})},
			Sphere{mgl32.Vec3{1, 0, -1}, 0.5, NewMetal(mgl32.Vec3{0.8, 0.6, 0.2}, 0.3)},
			Sphere{mgl32.Vec3{-1, 0, -1}, 0.5, NewDielectric(1.5)},
			Sphere{mgl32.Vec3{-1, 0, -1}, -0.45, NewDielectric(1.5)},
		},
	}

	img := image.NewRGBA(image.Rect(0, 0, nx, ny))
	for j := ny - 1; j >= 0; j-- {
		for i := 0; i < nx; i++ {
			col := mgl32.Vec3{0, 0, 0}
			for s := 0; s < ns; s++ {
				fi, fj := float32(i), float32(j)
				u, v := (fi+rand.Float32())/float32(nx), (fj+rand.Float32())/float32(ny)
				ray := cam.GetRay(u, v)
				col = Vadd(col, CalcColor(ray, world, 0))
			}
			col = Vdiv(col, float32(ns))
			col = mgl32.Vec3{Sqrt32(col.X()), Sqrt32(col.Y()), Sqrt32(col.Z())} // gamma 2.0
			img.Set(i, ny-j, ConvToColor(col))                                  // reverse height
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
