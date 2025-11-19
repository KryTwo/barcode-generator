package layout

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"main/config"
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

// when saveToFile == true, the function returns nil and saves the .pdf file.
func MakePDF(img []image.Image, data [][]string, saveToFile bool) []byte {
	var xPosTemp, yPosTemp float64

	cfg := config.Get()
	//требуемые параметры баркода и ячейки
	higth := cfg.Hight               //мм
	width := cfg.Width               //мм
	ySpacing := cfg.YSpacing         //pt
	xSpacing := cfg.XSpacing         //pt
	margin := float64(cfg.Margin)    //pt
	cellSizeMultiplier := 1.1        //множитель размера белого фона
	marginToCrop := cfg.MarginToCrop //отступ в бок от краев штрихкода для прорисовки линии нарезки листа

	//размеры баркода в мм
	bcHight := convert.MMToPointPDF(higth)
	bcWidth := convert.MMToPointPDF(width)

	//cfg := config.Get()
	pdf := gofpdf.New("p", "pt", "A4", "")

	pdf.AddPage()

	//загружаем шрифт из .json и .z
	loadFont(pdf)
	pdf.SetFont("DejaVuSans", "", float64(cfg.FontSize))

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
			//рисуем маркеры для резки
			xPosToCrop := pdf.GetX() - float64(cfg.MarginToCrop)
			yPosToCrop := pdf.GetY()
			pdf.SetLineWidth(0.2)
			//левый верхний маркер
			pdf.Line(xPosToCrop, yPosToCrop, xPosToCrop+float64(marginToCrop), yPosToCrop)
			pdf.Line(xPosToCrop, yPosToCrop, xPosToCrop, yPosToCrop+float64(marginToCrop))
			//правый нижний маркер
			pdf.Line(xPosToCrop+bcWidth+float64(marginToCrop), yPosToCrop+bcHight, xPosToCrop+bcWidth+float64(marginToCrop)*2, yPosToCrop+bcHight)
			pdf.Line(xPosToCrop+bcWidth+float64(marginToCrop)*2, yPosToCrop+bcHight, xPosToCrop+bcWidth+float64(marginToCrop)*2, yPosToCrop+bcHight-float64(marginToCrop))
			// fmt.Printf("data: %v\n", data)
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

			//сохраняем текущие координаты
			xPosTemp, yPosTemp = pdf.GetXY()

			//добавляем фон
			pdf.SetFillColor(255, 255, 255)
			textWidht := pdf.GetStringWidth(data[i][1])
			textHight, _ := pdf.GetFontSize()
			textHight = textHight * cellSizeMultiplier
			textWidht = textWidht * cellSizeMultiplier

			//размещаем текст
			// pdf.SetX(margin + bcWidth/2 - textWidht/2)
			pdf.SetX(xPos + bcWidth/2 - textWidht/2)
			pdf.CellFormat(textWidht, textHight, tr(data[i][1]), "", 0, "C", true, 0, "")

			//возвращаем координаты исходной точки
			pdf.SetXY(xPosTemp, yPosTemp)

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

	switch saveToFile {
	case false:
		var buf bytes.Buffer
		pdf.Output(&buf)
		pdfBytes := buf.Bytes()
		return pdfBytes
	default:
		err := pdf.OutputFileAndClose("resultToPrint.pdf")
		if err != nil {
			fmt.Printf("outpuFileAndClose error: %v\n", err)
		}
		return nil
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
