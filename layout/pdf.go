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
	//размеры баркода в мм
	bcHight := convert.MMToPointPDF(30)
	bcWidth := convert.MMToPointPDF(70)

	//cfg := config.Get()
	pdf := gofpdf.New("p", "pt", "A4", "")
	type countryType struct {
		nameStr, capitalStr, areaStr, popStr string
	}
	pdf.AddPage()

	//стартовая точка
	pdf.SetXY(25, 25)

	//размеры ячейки
	xCell := float64(convert.MMToPointPDF(70))
	yCell := float64(convert.MMToPointPDF(30))

	fmt.Println("pagesize w & h")
	fmt.Println(pdf.PageSize(1))

	xPos, yPos := pdf.GetXY()

	pdf.SetMargins(25, 25, 25)
	pdf.SetAutoPageBreak(false, 25)

	fmt.Println(pdf.GetMargins())

	improvedTable := func() {
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

			pdf.CellFormat(xCell, yCell, "", "1", 0, "L", false, 0, "")
			pdf.Ln(yCell + 10)

			pdf.RegisterImageOptionsReader(fileName, opt, strings.NewReader(imgBuf.String()))
			pdf.Image(fileName, xPos, yPos, bcWidth, bcHight, false, "", 0, "")
			yPos = pdf.GetY()
		}
	}
	pdf.SetFont("Arial", "", 18)
	improvedTable()

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
	// 	//отступ от границ листа
	// 	xBound := convert.MMToPT(15)
	// 	yBound := convert.MMToPT(15)

	// 	//стартовая точка
	// 	xPos := xBound
	// 	yPos := yBound

	// 	//отступ между штрихкодами
	// 	bcYSpace := convert.MMToPT(5)
	// 	bcXSpace := convert.MMToPT(15)

	// 	//TODO: акинуть текст с фоновой заливкой в ШК

	//закидываем баркоды на лист
	// for i := 0; i < len(img); i++ {
	// 	fileName := "barcode" + strconv.Itoa(i)
	// 	imgBuf, err := imageToPNG(img[i])
	// 	if err != nil {
	// 		fmt.Printf("err: %v\n", err)
	// 	}

	// 	opt := gofpdf.ImageOptions{
	// 		ImageType: "PNG",
	// 		ReadDpi:   true,
	// 	}

	// 	pdf.RegisterImageOptionsReader(fileName, opt, strings.NewReader(imgBuf.String()))
	// 	pdf.ImageOptions(fileName, float64(xPos), float64(yPos), bcWidth, bcHight, false, opt, 0, "")
	// 	yPos = yPos + bcYSpace + bcHight

	// 	//переренос на следующую колонку
	// 	if yPos+yBound+bcHight > 842 {
	// 		yPos = yBound
	// 		xPos = xPos + bcWidth + bcXSpace
	// 	}
	// 	//перенос на следующий лист
	// 	if xPos+bcWidth+bcXSpace > 595 {
	// 		fmt.Printf("граница достигнута на элементе i: %v\n", i+1)
	// 		pdf.AddPage()
	// 		xPos = xBound
	// 		yPos = yBound
	// 	}
	// 	}

	// }

}
