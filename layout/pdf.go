package layout

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"main/convert"
	"strconv"
	"strings"

	"github.com/jung-kurt/gofpdf"
)

// const mmInch = 25.4

// var dpi = 72

// преобразует миллиметры в пиксели, для размещения на листе.
// размер листа а4:
// при dpi:72  | (595 / 842)
// при dpi:96  | (794 / 1123)
// при dpi:300 | (2450 / 3508)
// при dpi:600 | (4960 / 7016)

func MakePDF(img []image.Image, data []string, maxX []int) {
	pdf := gofpdf.New("p", "pt", "A4", "")
	pdf.AddPage()

	//отступ от границ листа
	xBound := convert.MMToPT(15)
	yBound := convert.MMToPT(15)

	//стартовая точка
	xPos := xBound
	yPos := yBound

	//размеры баркода в мм
	bcHight := convert.MMToPT(30)
	bcWidth := convert.MMToPT(70)

	//отступ между штрихкодами
	bcYSpace := convert.MMToPT(5)
	bcXSpace := convert.MMToPT(15)

	//TODO: акинуть текст с фоновой заливкой в ШК

	//закидываем баркоды на лист
	for i := 0; i < len(img); i++ {
		fileName := "barcode" + strconv.Itoa(i)
		imgBuf, err := imageToPNG(img[i])
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}

		opt := gofpdf.ImageOptions{
			ImageType: "PNG",
			ReadDpi:   true,
		}

		pdf.RegisterImageOptionsReader(fileName, opt, strings.NewReader(imgBuf.String()))
		pdf.ImageOptions(fileName, float64(xPos), float64(yPos), bcWidth, bcHight, false, opt, 0, "")
		yPos = yPos + bcYSpace + bcHight

		//переренос на следующую колонку
		if yPos+yBound+bcHight > 842 {
			yPos = yBound
			xPos = xPos + bcWidth + bcXSpace
		}
		//перенос на следующий лист
		if xPos+bcWidth+bcXSpace > 595 {
			fmt.Printf("граница достигнута на элементе i: %v\n", i+1)
			pdf.AddPage()
			xPos = xBound
			yPos = yBound
		}
	}

	err := pdf.OutputFileAndClose("hello.pdf")
	if err != nil {
		fmt.Printf("outpuFileAndClose error: %v\n", err)
	}
}

// прогоняем image в буфер
func imageToPNG(img image.Image) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		return nil, err
	}
	return &buf, nil
}
