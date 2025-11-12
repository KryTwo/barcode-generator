package main

import (
	"main/app"
	"main/config"
	"main/gui"

	"fyne.io/fyne/v2"
	fyneApp "fyne.io/fyne/v2/app"
)

func main() {
	config.Init()
	cfg := config.Get()

	myApp := fyneApp.NewWithID("example.myapp")
	window := myApp.NewWindow("Barcode Generator")
	controller := app.NewController(cfg)
	gui.MakeUI(window, controller)

	window.Resize(fyne.Size{Width: 800, Height: 800})
	window.CenterOnScreen()
	window.ShowAndRun()

	//cfg := config.Get()
	//	filePath := "source/code.csv"

	//построчное получение данных
	// records, _, err := csvreader.Read(filePath)
	// if err != nil {
	// 	fmt.Println("ошибка чтения файла: ", err)
	// }
	// var arrRgba []image.Image
	// //генерация баркода
	// //срез максимальных ширин штрихкодов
	// var maximumX []int
	// var data []string
	// for i := 0; i < len(records); i++ {

	// 	img, err := barcode.GenerateCode128(records[i][0])
	// 	label.MakeFile(img, records[i][0])

	// 	bcLenX := img.Bounds().Max.X

	// 	maximumX = append(maximumX, bcLenX)
	// 	if err != nil {
	// 		fmt.Printf("can't generate code 128 with error: %v\n", err)
	// 	}

	// 	fmt.Println(records[i])
	// 	// rgba = label.DrawText(records[i][0], img, bcLenX)
	// 	arrRgba = append(arrRgba, img)
	// 	data = append(data, records[i][1])
	// }

	// //создаем файл PDF
	// layout.MakePDF(arrRgba, data, maximumX)

}
