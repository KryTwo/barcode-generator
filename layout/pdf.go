package layout

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"strings"

	"github.com/jung-kurt/gofpdf"
)

func MakePDF(img image.Image) {
	pdf := gofpdf.New("p", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	//pdf.Cell(40, 10, "hello world")
	imgBuf, _ := imageToPNG(img)
	// pdf.Image("out.png", 15, 0, 74, 30, true, "", 0, "")

	opt := gofpdf.ImageOptions{
		ImageType: "PNG",
		ReadDpi:   true,
	}
	pdf.RegisterImageOptionsReader("barcode", opt, strings.NewReader(imgBuf.String()))
	pdf.ImageOptions("barcode", 10, 20, 50, 0, false, opt, 0, "")

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
