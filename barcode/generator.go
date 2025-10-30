package barcode

import (
	"image"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
)

func GenerateCode128(data string, hight int, width int) (image.Image, error) {
	// bc := "CEL 3747872"

	bCode, err := code128.Encode(data)
	if err != nil {
		return nil, err
	}
	scaledBC, err := barcode.Scale(bCode, width, hight)
	if err != nil {
		return nil, err
	}
	return scaledBC, nil
	// BCfile, _ := os.Create("BCode.png")
	// defer BCfile.Close()
	// png.Encode(BCfile, scaledBC)

}
