package app

import (
	"fmt"
	"image"
	"main/barcode"
	"main/convert"
	"main/csvreader"
	"main/layout"
	"main/structs"
)

type Controller struct {
	config *structs.Config
}

type ProcessResult struct {
	PreviewPNG *image.RGBA
	PreviewBC  *image.RGBA
	Error      error
}

func NewController(config *structs.Config) *Controller {
	return &Controller{config: config}
}

func (c *Controller) ProcessFile(data []byte) ProcessResult {
	records, _, err := csvreader.Read(data)
	if err != nil {
		return ProcessResult{Error: err}
	}
	imgs, _ := barcode.GenerateCode128(records)

	PDFBytes := layout.MakePDF(imgs, records)
	previewBC := layout.BytesPdfToPNGConvert(PDFBytes)
	previewIMG := layout.BytesPdfToPNGConvert(PDFBytes)

	return ProcessResult{PreviewBC: previewBC, PreviewPNG: previewIMG}
}

func (c *Controller) CropBC(img *image.RGBA) *image.Image {
	fmt.Printf("img.Bounds(): %v\n", img.Bounds())
	fmt.Printf("c.config.Margin: %v\n", c.config.Margin)
	var x1, x2, y1, y2 float64
	x1 = c.config.Margin/72*c.config.DPI - 20
	y1 = c.config.Margin/72*c.config.DPI - 20
	x2 = x1 + convert.MMToPT(c.config.Width) + 40
	y2 = y1 + convert.MMToPT(c.config.Height) + 40
	croppRect := image.Rect(int(x1), int(y1), int(x2), int(y2))
	croppImg := img.SubImage(croppRect)

	return &croppImg
}
