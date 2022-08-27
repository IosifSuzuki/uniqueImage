package main

import (
	"fmt"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	color2 "image/color"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode"
	prng "uniqueImage/prng"
)

func main() {
	var userName = "iosifSuzuki"
	var image = generateUniqueImageBy(userName)
	targetFile := filepath.Join("output", fmt.Sprintf("%s.jpg", userName))
	file, err := os.Create(targetFile)
	defer func() {
		_ = file.Close()
	}()
	if err != nil {
		log.Fatal(err)
	}
	err = jpeg.Encode(file, image, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func key(userName string) int {
	var result = 0
	for _, a := range userName {
		result += int(a)
	}
	return result
}

func abbreviation(title string) string {
	if len(title) == 0 {
		return ""
	}
	var lastUppercaseLetterIndex = 0
	var abb = string(title[0])
	for idx, letter := range title {
		if (lastUppercaseLetterIndex+1 != idx || idx != 0) && unicode.IsUpper(letter) {
			abb += string(letter)
			lastUppercaseLetterIndex = idx
		}
	}
	return strings.ToUpper(abb)
}

func generateUniqueImageBy(userName string) *image.RGBA {
	const (
		width    = 100
		height   = 100
		fontSize = 50
	)
	var rectangle = image.Rect(0, 0, width, height)
	var imageTarget = image.NewRGBA(rectangle)
	var prng = prng.PRNG{
		Seed: key(userName),
	}
	var tmpColor = color2.RGBA{
		R: uint8(prng.GenerateNum() % 255),
		G: uint8(prng.GenerateNum() % 255),
		B: uint8(prng.GenerateNum() % 255),
		A: 255,
	}
	for row := 0; row < height; row++ {
		for column := 0; column < width; column++ {
			var color = color2.RGBA{
				R: uint8(prng.GenerateNum() % 255),
				G: uint8(prng.GenerateNum() % 255),
				B: uint8(prng.GenerateNum() % 255),
				A: 255,
			}
			var imageColor = color2.RGBA{
				R: tmpColor.R + (tmpColor.R-color.R)/2,
				G: tmpColor.G + (tmpColor.G-color.G)/2,
				B: tmpColor.B + (tmpColor.B-color.B)/2,
			}
			imageTarget.Set(row, column, imageColor)
		}
	}

	var fontFile = "source/Anton-Regular.ttf"
	fontBytes, err := ioutil.ReadFile(fontFile)
	if err != nil {
		log.Fatal(err)
	}
	f, err := truetype.Parse(fontBytes)
	if err != nil {
		log.Fatal(err)
	}

	drawer := font.Drawer{
		Dst: imageTarget,
		Src: image.NewUniform(color2.Black),
		Face: truetype.NewFace(f, &truetype.Options{
			Size:    fontSize,
			DPI:     0,
			Hinting: 0,
		}),
	}
	sizeOfText := drawer.MeasureString(abbreviation(userName))
	drawer.Dot = fixed.P((width-sizeOfText.Round())/2, (height+fontSize)/2)
	drawer.DrawString(abbreviation(userName))
	return imageTarget
}
