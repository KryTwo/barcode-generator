package layout

import (
	"barcode-app/logger"
	"image"

	"github.com/gen2brain/go-fitz"
)

func PdfToPNGConvert() *image.RGBA {
	doc, err := fitz.New("hello.pdf")
	if err != nil {
		logger.LogError(err, "error in pdfToPNGConvert (fitz.New)")

	}

	defer doc.Close()

	img, err := doc.Image(0)
	if err != nil {
		logger.LogError(err, "error in pdfToPNGConvert (dpc.Image)")
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
