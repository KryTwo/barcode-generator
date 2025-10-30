package barcode

import (
	"image"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
)

// генерирует barcode из data string с заданной шириной и высотой.
func GenerateCode128(data string, hight int, width int) (image.Image, error) {
	color := barcode.ColorScheme8
	bCode, err := code128.EncodeWithColor(data, color)
	if err != nil {
		return nil, err
	}

	scaledBC, err := barcode.Scale(bCode, width, hight)
	if err != nil {
		return nil, err
	}

	return scaledBC, nil
}
