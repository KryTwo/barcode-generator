package main

import (
	"flag"
	"fmt"
	"main/barcode"
	"main/csvreader"
	"main/label"
	"main/structs"
)

func ParseConfig() structs.Config {
	dpi := flag.Float64("dpi", 72, "screen resolution in Dots Per Inch")
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
		fmt.Println("ошибка чтения факла: ", err)
	}

	//генерация баркода
	//for _, v := range records {

	img, err := barcode.GenerateCode128(records[0][0], 100, 300)
	if err != nil {
		fmt.Printf("can't generate code 128 with error: %v\n", err)
	}
	fmt.Println(records[0][0])
	label.DrawText(records[0][0], cfg, img)
	//}

}
