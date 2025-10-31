package barcode

import (
	"image"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
)

// генерирует barcode из data string с заданной шириной и высотой.
func GenerateCode128(data string, width int, hight int) (image.Image, error, int) {
	color := barcode.ColorScheme8
	bCode, err := code128.EncodeWithColor(data, color)
	if err != nil {
		return nil, err, 0
	}

	maxX := bCode.Bounds().Max.X

	scaledBC, err := barcode.Scale(bCode, maxX, hight)
	if err != nil {
		return nil, err, 0
	}

	return scaledBC, nil, maxX
}
