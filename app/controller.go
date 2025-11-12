package app

import (
	"image"
	"main/barcode"
	"main/csvreader"
	"main/layout"
	"main/structs"
)

type Controller struct {
	config *structs.Config
}

type ProcessResult struct {
	PreviewPNG *image.RGBA
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

	layout.MakePDF(imgs, records)

	previewIMG := layout.PdfToPNGConvert()

	return ProcessResult{previewIMG, err}
}
