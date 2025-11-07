package barcode

import (
	"image"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
)

// генерирует barcode из data string
func GenerateCode128(data string) (image.Image, error) {
	bCode, err := code128.EncodeWithColor(data, barcode.ColorScheme8)
	if err != nil {
		return nil, err
	}

	return bCode, nil
}
