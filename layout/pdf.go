package layout

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"main/convert"
	"os"
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
	//требуемые параметры баркода и ячейки
	higth := 20               //мм
	width := 50               //мм
	ySpacing := 30.0          //pt
	xSpacing := 30.0          //pt
	margin := 50.0            //pt
	fontSize := 16.0          //
	cellSizeMultiplier := 1.1 //множитель размера белого фона

	//размеры баркода в мм
	bcHight := convert.MMToPointPDF(higth)
	bcWidth := convert.MMToPointPDF(width)

	//cfg := config.Get()
	pdf := gofpdf.New("p", "pt", "A4", "")

	pdf.AddPage()

	//загружаем шрифт из .json и .z
	loadFont(pdf)
	pdf.SetFont("DejaVuSans", "", fontSize)

	//стартовая точка
	pdf.SetXY(margin, margin)
	tr := pdf.UnicodeTranslatorFromDescriptor("./fonts/cp1251")

	// 	pagesize w & h
	// 595.28 841.89 pt
	xPageSize, yPageSize := pdf.GetPageSize()

	xPos, yPos := pdf.GetXY()

	//отступ от границ листа
	pdf.SetMargins(margin, margin, margin)
	pdf.SetAutoPageBreak(false, margin)

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

			pdf.RegisterImageOptionsReader(fileName, opt, strings.NewReader(imgBuf.String()))
			pdf.Image(fileName, xPos, yPos, bcWidth, bcHight, false, "", 0, "")

			xPosTemp, yPosTemp := pdf.GetXY()
			pdf.SetFillColor(255, 255, 255)

			textWidht := pdf.GetStringWidth(data[i])
			textHight, _ := pdf.GetFontSize()
			textHight = textHight * cellSizeMultiplier
			textWidht = textWidht * cellSizeMultiplier

			// pdf.SetX(margin + bcWidth/2 - textWidht/2)
			pdf.SetX(xPos + bcWidth/2 - textWidht/2)
			pdf.CellFormat(textWidht, textHight, tr(data[i]), "", 0, "C", true, 0, "")

			//возвращаем координаты исходной точки
			pdf.SetY(yPosTemp)
			pdf.SetX(xPosTemp)

			pdf.Ln(bcHight + ySpacing)

			yPos = pdf.GetY()

			//смещение на второй столбец
			if yPos >= yPageSize-ySpacing-bcHight {
				fmt.Printf("выход за пределы по высоте, итерация: %v\n\n", i)
				pdf.SetY(margin)
				yPos = pdf.GetY()

				pdf.SetX(xPos + xSpacing + bcWidth)
				xPos = pdf.GetX()
			}

			//смещение в начало нового листа
			if xPos >= xPageSize-xSpacing-bcWidth {
				fmt.Printf("выход за пределы по ширине, итерация: %v\n\n", i)
				pdf.AddPage()
				pdf.SetXY(margin, margin)
				xPos = pdf.GetX()
				yPos = pdf.GetY()
			}
			// fmt.Printf("xPos: %v ", math.Round(pdf.GetX()))
			// fmt.Printf("yPos: %v\n", math.Round(yPos))
		}
	}

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
}

func loadFont(pdf *gofpdf.Fpdf) {
	jsonBytes, err := os.ReadFile("./fonts/DejaVuSans.json")
	if err != nil {
		fmt.Printf("jsonBytes err: %v\n", err)
	}
	zBytes, err := os.ReadFile("./fonts/DejaVuSans.z")
	if err != nil {
		fmt.Printf("zBytes err: %v\n", err)
	}

	pdf.AddFontFromBytes("DejaVuSans", "", jsonBytes, zBytes)
}
