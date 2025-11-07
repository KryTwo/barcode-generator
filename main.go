package main

import (
	"fmt"
	"image"
	"main/barcode"
	"main/config"
	"main/csvreader"
	"main/label"
	"main/layout"
)

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
	config.Init()
	//cfg := config.Get()
	filePath := "source/code.csv"

	//построчное получение данных

	records, _, err := csvreader.Read(filePath)
	if err != nil {
		fmt.Println("ошибка чтения файла: ", err)
	}
	// var rgba *image.RGBA
	var arrRgba []image.Image
	//генерация баркода
	//срез максимальных ширин штрихкодов
	var maximumX []int
	var data []string
	for i := 0; i < len(records); i++ {

		img, err := barcode.GenerateCode128(records[i][0])
		label.MakeFile(img, records[i][0])

		bcLenX := img.Bounds().Max.X

		maximumX = append(maximumX, bcLenX)
		if err != nil {
			fmt.Printf("can't generate code 128 with error: %v\n", err)
		}

		fmt.Println(records[i])
		// rgba = label.DrawText(records[i][0], img, bcLenX)
		arrRgba = append(arrRgba, img)
		data = append(data, records[i][0])
	}

	//создаем файл PDF
	layout.MakePDF(arrRgba, data, maximumX)

}
