package barcode

import (
	"fmt"
	"image"
	"main/config"
	"main/convert"
	"main/label"
	"math"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
)

// генерирует barcode из data string с заданной шириной и высотой.
func GenerateCode128(data string) (image.Image, error) {
	cfg := config.Get()

	bCode, err := code128.EncodeWithColor(data, barcode.ColorScheme8)
	if err != nil {
		return nil, err
	}

	//ширина баркода в Pt
	//bcLenX := bCode.Bounds().Max.X
	fmt.Printf("cfg.Width: %v mm or %v pt\n", cfg.Width, math.Round(convert.MMToPT(cfg.Width)))

	// Скейл баркода по высоте в пикселях
	// scaledBC, err := barcode.Scale(bCode, bcLenX, int(convert.MMToPT(cfg.Height)))
	// scaledBC, err := barcode.Scale(bCode, bCode.Bounds().Min.X, int(cfg.Height))
	// if err != nil {
	// 	return nil, err
	// }
	label.MakeFile(bCode, data)
	return bCode, nil
}
