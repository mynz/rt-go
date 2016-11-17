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

func RenderImage() *image.RGBA {
	nx, ny := 200, 100

	img := image.NewRGBA(image.Rect(0, 0, nx, ny))

	for j := ny-1; j >= 0; j-- {
		for  i := 0; i < nx; i++ {
			col := mgl32.Vec3{ float32(i) / float32(nx), float32(j) / float32(ny), 0.2 }
			c := col.Mul(255.99)
			color := color.RGBA{ uint8(c[0]), uint8(c[1]), uint8(c[2]), 255 }
			img.Set(i, j, color)
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


