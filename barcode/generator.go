package barcode

import (
	"image"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
)

// генерирует barcode из data string
func GenerateCode128(data [][]string) ([]image.Image, error) {
	var bCodeList []image.Image
	for _, s := range data {
		bCode, err := code128.EncodeWithColor(s[0], barcode.ColorScheme8)
		if err != nil {
			return nil, err
		}
		bCodeList = append(bCodeList, bCode)
	}
	return bCodeList, nil
}
