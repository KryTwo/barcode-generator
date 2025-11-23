package main

import (
	"main/app"
	"main/config"
	"main/gui"
	"main/logger"

	fyneApp "fyne.io/fyne/v2/app"
)

func main() {
	logger.Log.Info("start main")

	config.Init()
	cfg := config.Get()

	myApp := fyneApp.NewWithID("bcgen.myapp")
	window := myApp.NewWindow("Barcode Generator")
	controller := app.NewController(cfg)
	gui.MakeUI(window, controller)

	//window.Resize(fyne.Size{Width: 800, Height: 800})
	window.CenterOnScreen()
	window.ShowAndRun()

}
