package app

import (
	"fmt"
	"image"
	"log"
	"main/barcode"
	"main/config"
	"main/convert"
	"main/csvreader"
	"main/layout"
	"main/structs"
	"strconv"
)

type Controller struct {
	config           *structs.Config
	CurrentRecords   [][]string
	OnPreviewUpdated func(*image.RGBA)
}

type ProcessResult struct {
	PreviewPNG *image.RGBA
	PreviewBC  *image.RGBA
	Success    bool
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

	c.CurrentRecords = records

	c.RegeneratePreview()

	return ProcessResult{Success: true}
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

func (c *Controller) SetBCWidth(data string) {
	d, err := strconv.Atoi(data)
	if err != nil {
		log.Fatalf("Failed convert ATOI in SetBCWidth: %v\n", err)
	}
	config.SetWidth(d)
	c.RegeneratePreview()
}

func (c *Controller) RegeneratePreview() {
	if len(c.CurrentRecords) == 0 {
		return
	}

	imgs, err := barcode.GenerateCode128(c.CurrentRecords)
	if err != nil {
		log.Fatalf("err: %v\n", err)
		return
	}

	PDFBytes := layout.MakePDF(imgs, c.CurrentRecords)

	img := layout.BytesPdfToPNGConvert(PDFBytes)

	if c.OnPreviewUpdated != nil {
		c.OnPreviewUpdated(img)
	}
}
