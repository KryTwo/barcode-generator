package main

import (
	"flag"
	"fmt"
	"image"
	"main/barcode"
	"main/csvreader"
	"main/label"
	"main/layout"
	"main/structs"
)

func ParseConfig() structs.Config {
	dpi := flag.Float64("dpi", 48, "screen resolution in Dots Per Inch")
	fontfile := flag.String("fontfile", "./fonts/RobotoforLearning-Black_0.ttf", "filename of the ttf font")
	hinting := flag.String("hinting", "none", "none | full")
	size := flag.Float64("size", 32, "font size in points")
	spacing := flag.Float64("spacing", 1.5, "line spacing (e.g. 2 means double spaced)")
	wonb := flag.Bool("whiteonblack", false, "white text on a black background")

	flag.Parse()
	return structs.Config{
		DPI:      *dpi,
		FontFile: *fontfile,
		Hinting:  *hinting,
		Size:     *size,
		Spacing:  *spacing,
		WONB:     *wonb,
	}
}

/*

	func main() {
		a := app.New()
		w := a.NewWindow("Hello World")
		w2 := a.NewWindow("window 2")

		message := widget.NewLabel("Welcome")
		message2 := widget.NewLabel("Welcome Twice")
		button := widget.NewButton("Update", func() {
			formatted := time.Now().Format("Time: 03:04:05")
			message.SetText(formatted)
		})

		w.SetContent(container.NewVBox(message, button))
		w.SetMaster()
		w.Resize(fyne.NewSize(350, 350))
		w.Show()

		buttonNewWindow := widget.NewButton("new Window", func() {
			w3 := a.NewWindow("new window")
			w3.SetContent(message2)
			w3.Show()
		})
		w2.SetContent(container.NewHBox(buttonNewWindow))
		w2.Resize(fyne.NewSize(350, 350))
		w2.Show()

		a.Run()
	}
*/

func main() {
	cfg := ParseConfig()

	//построчное получение данных

	records, _, err := csvreader.Read("source/code.csv")
	if err != nil {
		fmt.Println("ошибка чтения файла: ", err)
	}
	var rgba *image.RGBA
	var arrRgba []image.Image
	//генерация баркода
	//срез максимальных ширин штрихкодов
	var maximumX []int
	for i := 0; i < 2; i++ {

		img, err, maxX := barcode.GenerateCode128(records[i][0], 300, 100)
		// label.MakeFile(img)
		maximumX = append(maximumX, maxX)
		if err != nil {
			fmt.Printf("can't generate code 128 with error: %v\n", err)
		}
		fmt.Println(records[i][0])
		rgba = label.DrawText(records[i][0], cfg, img, maxX)
		arrRgba = append(arrRgba, rgba)
	}

	//создаем файл PDF
	layout.MakePDF(arrRgba, maximumX)

}
