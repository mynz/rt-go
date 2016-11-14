package main

import(
	"os"
	"fmt"
	"image"
	"image/color"
	"image/png"
)

func writePngFile(filename string) {
	w, h := 128, 128
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	if true {
		c := color.RGBA{ 255, 128, 0, 255 }
		for i := 0; i < w; i++ {
			img.Set(i, i, c)
		}
	}

	fp, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	if err := png.Encode(fp, img); err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println("Hello.")

	writePngFile("test.png")

	fmt.Println("Done.")
}


