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

func writePngFile(filename string, img image.Image) {

/*
 *     w, h := 128, 128
 *     img := image.NewRGBA(image.Rect(0, 0, w, h))
 * 
 *     if true {
 *         c := color.RGBA{ 255, 128, 0, 255 }
 *         for i := 0; i < w; i++ {
 *             img.Set(i, i, c)
 *         }
 *     }
 */

	fp, err := os.Create(filename)
	defer fp.Close()
	if err != nil {
		panic(err)
	}

	if err := png.Encode(fp, img); err != nil {
		panic(err)
	}
}

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


func main() {
	fmt.Println("Hello.")

	img := RenderImage()
	writePngFile("test.png", img)

/*
 *     if false {
 *         // fmt.Println("Vec:", vecmath.NewVecZero().Length())
 * 
 *         v := mgl32.Vec3{1, 2, 3}
 *         fmt.Println("v0:", v.Len())
 * 
 *         fmt.Println("v1:", mgl32.Vec3{4, 5, 6}.Len())
 * 
 *         writePngFile("test.png")
 *     }
 */

	fmt.Println("Done.")
}


