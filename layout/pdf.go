package layout

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"strconv"
	"strings"

	"github.com/jung-kurt/gofpdf"
)

const mmInch = 25.4

func mmToPt(mm int, dpi int) float64 {
	var pt float64
	pt = float64(mm) / mmInch * 72.0
	return pt
}

func MakePDF(img []image.Image, maxX []int) {
	pdf := gofpdf.New("p", "pt", "A4", "")
	pdf.AddPage()

	xPos := 50
	yPos := 0
	dpi := 72
	bcWidth := mmToPt(70, dpi)

	for i := 0; i < len(img); i++ {
		fileName := "barcode" + strconv.Itoa(i)
		imgBuf, _ := imageToPNG(img[i])

		opt := gofpdf.ImageOptions{
			ImageType:             "PNG",
			ReadDpi:               true,
			AllowNegativePosition: false,
		}

		pdf.RegisterImageOptionsReader(fileName, opt, strings.NewReader(imgBuf.String()))
		pdf.ImageOptions(fileName, float64(xPos), float64(yPos), bcWidth, -1, false, opt, 0, "")
		yPos += 100
	}
	err := pdf.OutputFileAndClose("hello.pdf")
	if err != nil {
		fmt.Printf("outpuFileAndClose error: %v\n", err)
	}
}

func imageToPNG(img image.Image) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		return nil, err
	}
	return &buf, nil
}
