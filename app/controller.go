package app

import (
	"main/barcode"
	"main/csvreader"
	"main/layout"
	"main/structs"
)

type Controller struct {
	config *structs.Config
}

func NewController(config *structs.Config) *Controller {
	return &Controller{config: config}
}

func (c *Controller) ProcessFile(data []byte) error {
	records, _, err := csvreader.Read(data)
	if err != nil {
		return err
	}
	imgs, _ := barcode.GenerateCode128(records)

	layout.MakePDF(imgs, records)

	layout.PdfToPNGConvert()

	return err
}
