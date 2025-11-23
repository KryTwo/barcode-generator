package layout

import (
	"image"
	"main/logger"

	"github.com/gen2brain/go-fitz"
)

func PdfToPNGConvert() *image.RGBA {
	doc, err := fitz.New("hello.pdf")
	if err != nil {
		panic(err)
	}

	defer doc.Close()

	img, err := doc.Image(0)
	if err != nil {
		panic(err)
	}

	return img
}

func BytesPdfToPNGConvert(b []byte) *image.RGBA {
	doc, err := fitz.NewFromMemory(b)
	if err != nil {
		logger.LogError(err, "err fitz.NewFromMemory")
	}

	defer doc.Close()

	img, err := doc.Image(0)
	if err != nil {
		logger.LogError(err, "err doc.Image")
	}

	return img
}
