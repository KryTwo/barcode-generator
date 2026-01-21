package main

import (
	"barcode-app/app"
	"barcode-app/config"
	"barcode-app/gui"
	"barcode-app/logger"

	fyneApp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
)

func main() {
	logger.Log.Info("start main")

	config.Init()
	cfg := config.Get()

	myApp := fyneApp.NewWithID("bcgen.myapp")
	myApp.Settings().SetTheme(theme.DarkTheme())
	window := myApp.NewWindow("Barcode Generator")
	controller := app.NewController(cfg)
	gui.MakeUI(window, controller, myApp)

	//window.Resize(fyne.Size{Width: 800, Height: 800})
	window.CenterOnScreen()
	window.ShowAndRun()

}
