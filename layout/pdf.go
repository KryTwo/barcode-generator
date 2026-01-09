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
	higth := cfg.Higth                        //мм
	width := cfg.Width                        //мм
	ySpacing := cfg.YSpacing                  //pt
	xSpacing := cfg.XSpacing                  //pt
	margin := float64(cfg.Margin)             //pt
	marginToCrop := float64(cfg.MarginToCrop) //отступ в бок от краев штрихкода для прорисовки линии нарезки листа
	originalFontSize := cfg.FontSize

	//размеры баркода в мм
	bcHigth := convert.MMToPointPDF(higth)
	bcWidth := convert.MMToPointPDF(width) - float64(marginToCrop)*2

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
			fmt.Printf("data[i][1]: %v\n", data[i][1])
			currentFontSize := originalFontSize
			pdf.SetFontSize(float64(currentFontSize))

			//рисуем маркеры для резки
			xPosToCrop := xPos - float64(cfg.MarginToCrop)
			yPosToCrop := yPos
			//толщина линий нарезки
			pdf.SetLineWidth(0.2)
			//левый верхний маркер
			pdf.Line(xPosToCrop, yPosToCrop, xPosToCrop+marginToCrop, yPosToCrop)
			pdf.Line(xPosToCrop, yPosToCrop, xPosToCrop, yPosToCrop+marginToCrop)
			//правый нижний маркер
			pdf.Line(xPosToCrop+bcWidth+marginToCrop, yPosToCrop+bcHigth, xPosToCrop+bcWidth+marginToCrop*2, yPosToCrop+bcHigth)
			pdf.Line(xPosToCrop+bcWidth+marginToCrop*2, yPosToCrop+bcHigth, xPosToCrop+bcWidth+marginToCrop*2, yPosToCrop+bcHigth-marginToCrop)

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
			pdf.Image(fileName, xPos, yPos, bcWidth, bcHigth, false, "", 0, "")

			//сохраняем текущие координаты
			xPosTemp, yPosTemp = pdf.GetXY()
			fmt.Printf("ширина текста: %v\n", pdf.GetStringWidth(data[i][1]))
			fmt.Printf("ширина баркода (после границ нарезки): %v\n", bcWidth)
			fmt.Println()
			//рисуем текст поверх шк
			drawBarcodeTextNew(pdf, tr, data[i][1], xPos, yPos, bcWidth, bcHigth)

			//возвращаем координаты исходной точки
			pdf.SetXY(xPosTemp, yPosTemp)

			//смещение координат для начала отрисовки следующего штрихкода
			pdf.Ln(bcHigth + ySpacing)

			yPos = pdf.GetY()

			//смещение на второй столбец, если текущий заполнен
			if yPos >= yPageSize-ySpacing-bcHigth {
				fmt.Printf("выход за пределы по высоте, итерация: %v\n\n", i)
				pdf.SetY(margin)
				yPos = pdf.GetY()

				pdf.SetX(xPos + xSpacing + bcWidth)
				xPos = pdf.GetX()
			}

			//смещение в начало нового листа, если текущий заполнен
			if xPos >= xPageSize-xSpacing-bcWidth {
				fmt.Printf("выход за пределы по ширине, итерация: %v\n\n", i)
				pdf.AddPage()
				pdf.SetXY(margin, margin)
				xPos = pdf.GetX()
				yPos = pdf.GetY()
			}
		}
	}

	improvedTable()

	if saveToFile {
		err := pdf.OutputFileAndClose("resultToPrint.pdf")
		if err != nil {
			fmt.Printf("outpuFileAndClose error: %v\n", err)
		}
		return nil
	} else {
		var buf bytes.Buffer
		pdf.Output(&buf)
		pdfBytes := buf.Bytes()
		return pdfBytes
	}

}

func drawBarcodeTextNew(pdf *gofpdf.Fpdf, tr func(string) string, text string, x, y, bcWidth, bcHigth float64) {
	cfg := config.Get()
	currentFontSize := cfg.FontSize
	//cellSizeMultiplier := 1.0
	textWidth := pdf.GetStringWidth(text)
	textHigth, _ := pdf.GetFontSize()
	pdf.SetFillColor(0, 255, 0)

	//Размещение текста с переносом по сепаратору
	if cfg.TextWrapping {
		//размещаем в текущих размерах, если не выходит за границы
		if textWidth < bcWidth && textHigth < bcHigth {
			pdf.SetX(x + (bcWidth-textWidth)/2)
			pdf.CellFormat(textWidth, textHigth, tr(text), "", 0, "C", true, 0, "")
		} else {
			//определяем разделитель
			sep := currentSeparator(text)

			//разбиваем текст на части по разделителю sep
			dataParts := strings.Split(text, sep)

			//добавляем разделитель для каждой из частей текста
			for i := range dataParts {
				if i != len(dataParts)-1 {
					dataParts[i] = dataParts[i] + sep
				}
			}

			//разбиваем текст на строки
			var buf string
			var result []string
			for {
				result = []string{}
				buf = ""
				textHigth, _ = pdf.GetFontSize()
				tooWide := false

				pdf.SetFontSize(float64(currentFontSize))

				//text = 132456789-5555-987654321
				for _, v := range dataParts {
					if pdf.GetStringWidth(v) > bcWidth {
						tooWide = true
						break
					}

					if pdf.GetStringWidth(buf+v) <= bcWidth {
						buf = buf + v
					} else {
						result = append(result, buf)
						buf = v
					}
				}
				if !tooWide {
					result = append(result, buf)

					totalHeight := len(result) * int(textHigth)

					if totalHeight <= int(bcHigth)/2 {
						break
					}
				}
				currentFontSize -= 1
				if currentFontSize < 1 {
					break
				}

			}

			//находим самую широкую строку
			var maxWidthString string
			for _, v := range result {
				if len(v) > len(maxWidthString) {
					maxWidthString = v
				}
			}

			textWidth = pdf.GetStringWidth(maxWidthString)
			fmt.Printf("textWidth: %v\n", textWidth)
			fmt.Printf("bcWidth: %v\n", bcWidth)
			for _, v := range result {
				pdf.SetX(x + (bcWidth-textWidth)/2)
				pdf.CellFormat(textWidth, textHigth, tr(v), "", 1, "C", true, 0, "")
			}
		}
	} else {
		//Готово
		//Размещение текста без переноса, с уменьшением размера текста
		if textWidth > bcWidth || textHigth > bcHigth/2 {
			for {
				currentFontSize -= 1
				pdf.SetFontSize(float64(currentFontSize))
				// fmt.Printf("currentFontSize: %v\n", currentFontSize)
				textWidth = pdf.GetStringWidth(text)
				textHigth, _ = pdf.GetFontSize()

				if textWidth < bcWidth && textHigth < bcHigth/2 {
					break
				}
			}
			pdf.SetX(x + (bcWidth-textWidth)/2)
			pdf.CellFormat(textWidth, textHigth, tr(text), "", 0, "C", true, 0, "")
		} else {
			pdf.SetX(x + (bcWidth-textWidth)/2)
			pdf.CellFormat(textWidth, textHigth, tr(text), "", 0, "C", true, 0, "")
		}
	}

}

func currentSeparator(text string) string {
	var sepList = []string{" ", ",", ".", "-"}
	var sepSymbol string
	for _, v := range sepList {
		out := 0
		count := strings.Count(text, v)
		if count > out {
			out = count
			sepSymbol = v
		}
	}
	return sepSymbol
}

/*func getRuneCount(data string) (int, int, int, int, int) {
	var latinUpper, latinLower, cyrillicUpper, cyrillicLower int
	var count int
	for _, r := range data {
		count++
		if unicode.Is(unicode.Latin, r) {
			if unicode.IsUpper(r) {
				latinUpper++
			} else {
				latinLower++
			}
		} else if unicode.Is(unicode.Cyrillic, r) {
			if unicode.IsUpper(r) {
				cyrillicUpper++
			} else {
				cyrillicLower++
			}
		}
	}
	fmt.Printf("data: %v\n", data)
	fmt.Printf("cyrillicLower: %v\n", cyrillicLower)
	fmt.Printf("cyrillicUpper: %v\n", cyrillicUpper)
	fmt.Printf("latinLower: %v\n", latinLower)
	fmt.Printf("latinUpper: %v\n", latinUpper)
	fmt.Printf("count: %v\n", count)
	fmt.Println()

	return latinUpper, latinLower, cyrillicUpper, cyrillicLower, count
}
*/
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
