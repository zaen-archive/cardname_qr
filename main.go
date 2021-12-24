package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
)

func getCard(filepath string) *image.RGBA {
	oriFile, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer oriFile.Close()

	oriPng, err := png.Decode(oriFile)
	if err != nil {
		panic(err)
	}

	oriBound := oriPng.Bounds()
	modPng := image.NewRGBA(image.Rect(0, 0, oriBound.Dx(), oriBound.Dy()))

	draw.Draw(modPng, oriBound, oriPng, oriBound.Min, draw.Src)

	return modPng
}

func main() {
	card := getCard("card.png")
	file, _ := os.Create("card2.png")

	defer file.Close()

	png.Encode(file, card)

	fmt.Println("~~ Finish ~~")
}
