package main

import (
	"fmt"
	"main/csvreader"
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

	records, _, err := csvreader.Read("source/code.csv")
	if err != nil {
		fmt.Println("ошибка чтения факла: ", err)
	}

	for v := range records {
		fmt.Println(v)
	}

	// bc := "CEL 3747872"
	// bCode, _ := code128.Encode(bc)
	// scaledBC, _ := barcode.Scale(bCode, 300, 100)
	// BCfile, _ := os.Create("BCode.png")
	// defer BCfile.Close()
	// png.Encode(BCfile, scaledBC)
}
