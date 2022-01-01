package main

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/lithammer/shortuuid/v3"
	"github.com/signintech/gopdf"
	"github.com/skip2/go-qrcode"
)

func createCardName(filepath string) *image.RGBA {
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

func addQrCode(card *image.RGBA) {
	id := shortuuid.New()
	url := os.Getenv("QR_PREFIX") + id

	x, _ := strconv.Atoi(os.Getenv("QR_X"))
	y, _ := strconv.Atoi(os.Getenv("QR_Y"))
	size, _ := strconv.Atoi(os.Getenv("QR_SIZE"))

	qrSource, err := qrcode.Encode(url, qrcode.Highest, size)
	if err != nil {
		panic(err)
	}

	qrImg, err := png.Decode(bytes.NewReader(qrSource))
	if err != nil {
		panic(err)
	}

	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			color := qrImg.At(i, j)
			card.Set(x+i, y+j, color)
		}
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	for i := 0; i < 10; i++ {
		card := createCardName(os.Getenv("CARD_PATH"))

		addQrCode(card)

		file, _ := os.Create("temp/card/card-" + strconv.Itoa(i) + ".png")

		png.Encode(file, card)
		file.Close()
	}

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()

	cardWidth, _ := strconv.ParseFloat(os.Getenv("CARD_WIDTH"), 64)
	cardHeight, _ := strconv.ParseFloat(os.Getenv("CARD_HEIGHT"), 64)
	// cardFlipped := os.Getenv("CARD_FLIP") == "true"

	cardWidth = gopdf.UnitsToPoints(gopdf.UnitMM, cardWidth)
	cardHeight = gopdf.UnitsToPoints(gopdf.UnitMM, cardHeight)
	pdfRect := &gopdf.Rect{
		W: cardWidth,
		H: cardHeight,
	}

	for i := 0; i < 2; i++ {
		for j := 0; j < 5; j++ {
			x := float64(i) * cardWidth
			y := float64(j) * cardHeight
			index := i*5 + j

			pdf.Image("temp/card/card-"+strconv.Itoa(index)+".png", x, y, pdfRect)
		}
	}

	pdf.WritePdf("temp/hello.pdf")
	fmt.Println("~~ Finish ~~")
}
