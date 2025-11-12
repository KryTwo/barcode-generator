package layout

import (
	"image"

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
