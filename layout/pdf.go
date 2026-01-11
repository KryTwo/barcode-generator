package layout

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"main/config"
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
	bcHigth := float64(higth)
	bcWidth := float64(width) - float64(marginToCrop)*2

	//cfg := config.Get()
	pdf := gofpdf.New("p", "mm", "A4", "")

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

			//визуальная линейка для упрощения (разделение по 1 см)
			xPosMarking := 10.0
			nums := []string{}

			for i := 1; i < 30; i++ {
				nums = append(nums, strconv.Itoa(i))
			}

			f := 0

			for xPosMarking <= xPageSize && f < len(nums) {
				pdf.Line(float64(xPosMarking), 1, float64(xPosMarking), 10)
				pdf.SetFontSize(14)
				pdf.SetTextColor(0, 0, 255)

				pdf.Text(xPosMarking-5, 5, nums[f])
				xPosMarking += 10

				f++
			}
			pdf.SetTextColor(0, 0, 0)
			///////////////////////////////////////////////

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
			//  ___ добавить недостающие маркеры (левый нижний и правый верхний)

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
			drawBarcodeText(pdf, tr, data[i][1], xPos, yPos, bcWidth, bcHigth)

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

	//если saveToFile == true, то сохраняем в файл, иначе копируем в буфер
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

func drawBarcodeText(pdf *gofpdf.Fpdf, tr func(string) string, text string, x, y, bcWidth, bcHigth float64) {
	cfg := config.Get()
	currentFontSize := cfg.FontSize
	//cellSizeMultiplier := 1.0
	textWidth := pdf.GetStringWidth(text)
	textHigth, _ := pdf.GetFontSize()
	pdf.SetFillColor(0, 255, 0)

	// fmt.Printf("textWidth: %v\n", textWidth)
	// fmt.Printf("bcWidth: %v\n", bcWidth)
	// fmt.Printf("textHigth: %v\n", textHigth)
	// fmt.Printf("bcHigth: %v\n", bcHigth)

	if cfg.TextWrapping {
		//Размещение текста с переносом по сепаратору
		if textWidth < bcWidth && textHigth < bcHigth {
			pdf.SetX(x + (bcWidth-textWidth)/2)
			pdf.CellFormat(textWidth, textHigth, tr(text), "", 0, "C", true, 0, "")
		} else {
			//находим разделитель
			sep := currentSeparator(text)

			//разбиваем текст на части по разделителю sep
			dataParts := strings.Split(text, sep)

			//добавляем разделитель для каждой из частей текста
			for i := range dataParts {
				if i != len(dataParts)-1 {
					dataParts[i] = dataParts[i] + sep
				}
			}
			var buf string
			var result []string
			var i int
			for {
				result = []string{}
				buf = ""
				textHigth, _ = pdf.GetFontSize()

				i++
				fmt.Println(i)

				pdf.SetFontSize(float64(currentFontSize))
				a, b := pdf.GetFontSize()

				fmt.Printf("a: %v\n", a)
				fmt.Printf("b: %v\n", b)

				//text = 123332221-5555-22
				for _, v := range dataParts {
					if pdf.GetStringWidth(buf+v) <= bcWidth {
						buf = buf + v
					} else {
						result = append(result, buf)
						buf = v
					}
				}
				result = append(result, buf)

				for _, v := range result {
					fmt.Printf("v: %v\n", v)
				}

				totalHeight := len(result) * int(textHigth)
				fmt.Printf("totalHeight: %v\n", totalHeight)
				fmt.Printf("bcHigth: %v\n", bcHigth)
				if totalHeight <= int(bcHigth)/2 {
					break
				} else {
					currentFontSize -= 1
				}

				fmt.Println()
			}

			//находим самую широкую строку
			var maxWidthString string
			for _, v := range result {
				if len(v) > len(maxWidthString) {
					maxWidthString = v
				}
			}

			textWidth = pdf.GetStringWidth(maxWidthString)

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

// func drawBarcodeText(pdf *gofpdf.Fpdf, tr func(string) string, text string, x, y, bcWidth, bcHigth float64) {
// 	cfg := config.Get()
// 	currentFontSize := cfg.FontSize
// 	cellSizeMultiplier := 1.0

// 	var dataParts []string
// 	//цвет фона
// 	pdf.SetFillColor(0, 255, 0)

// 	textWidth := pdf.GetStringWidth(text)
// 	textWidth = textWidth * cellSizeMultiplier

// 	textHigth, _ := pdf.GetFontSize()
// 	textHigth *= cellSizeMultiplier

// 	//var sepList = []string{" ", ",", ".", "-"}

// 	//переносим или уменьшаем шрифт
// 	if cfg.TextWrapping {
// 		//Перенос текста
// 		delimiter := "-" //TODO: добавить автоопределение разделителя
// 		if textWidth > float64(bcWidth) {
// 			dataParts = strings.Split(text, delimiter)
// 			pdf.SetFontSize(float64(currentFontSize))
// 		}
// 	} else {
// 		//уменьшаем размер текста
// 		for textWidth > float64(cfg.Width) || textHigth*2 > bcHigth {
// 			currentFontSize -= 1
// 			pdf.SetFontSize(float64(currentFontSize))
// 			textWidth = pdf.GetStringWidth(text)
// 			textHigth, _ = pdf.GetFontSize()
// 		}
// 	}

// 	//размещаем текст
// 	if len(dataParts) > 1 {
// 		for _, v := range dataParts {
// 			for len(dataParts)*int(textHigth) > int(bcHigth)/2 {
// 				textHigth, _ = pdf.GetFontSize()
// 				currentFontSize -= 1
// 				pdf.SetFontSize(float64(currentFontSize))
// 			}

// 			textWidth = pdf.GetStringWidth(v)

// 			pdf.SetX(float64(cfg.Margin) + bcWidth/2 - textWidth/2 - float64(cfg.MarginToCrop)/2)
// 			pdf.CellFormat(textWidth, textHigth, tr(v), "", 0, "C", true, 0, "")
// 			pdf.SetY(pdf.GetY() + textHigth)
// 		}

// 	} else {
// 		pdf.SetX(x + bcWidth/2 - textWidth/2)
// 		pdf.CellFormat(textWidth, textHigth, tr(text), "", 0, "C", true, 0, "")
// 	}
// }

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
