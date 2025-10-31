package barcode

import (
	"fmt"
	"image"
	"main/label"

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
	fmt.Printf("bCode.Bounds(): %v\n", bCode.Bounds().Max)
	label.MakeFile(bCode, ("before scaling" + data))
	maxX := bCode.Bounds().Max.X
	scaledBC, err := barcode.Scale(bCode, maxX, hight)
	if err != nil {
		return nil, err, 0
	}
	fmt.Printf("scaledBC.Bounds(): %v\n", scaledBC.Bounds().Max)
	label.MakeFile(scaledBC, data)

	return scaledBC, nil, maxX
}
