package main

import (
	"fmt"
	"image/png"
	"io/ioutil"
	"os"

	"github.com/danielgatis/go-findfont/findfont"
	"github.com/danielgatis/go-freetype/freetype"
)

func checkErr(err error) {
	if err != nil {
		fmt.Printf("Failed: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	fonts, err := findfont.Find("Arial", findfont.FontRegular)
	checkErr(err)

	data, err := ioutil.ReadFile(fonts[0][2])
	checkErr(err)

	lib, err := freetype.NewLibrary()
	checkErr(err)

	face, err := freetype.NewFace(lib, data, 0)
	checkErr(err)

	err = face.Pt(32, 72)
	checkErr(err)

	img, m, err := face.Glyph('A')
	checkErr(err)

	f, err := os.Create("image.png")
	checkErr(err)

	err = png.Encode(f, img)
	checkErr(err)

	fmt.Printf("Image: %v\n", f.Name())
	fmt.Printf("Metrics: %v\n", m)

	err = face.Done()
	checkErr(err)

	err = lib.Done()
	checkErr(err)
}
