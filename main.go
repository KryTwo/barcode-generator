package main

import (
	"barcode-app/app"
	"barcode-app/config"
	"barcode-app/gui"
	"barcode-app/logger"

	fyneApp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
)

var SettingsJSON config.JSONSettings

func main() {

	logger.Log.Info("start main")

	config.Init()
	cfg := config.Get()
	myApp := fyneApp.NewWithID("com.krytwo.vi-helper.bcvigen")
	myApp.Settings().SetTheme(theme.DarkTheme())
	window := myApp.NewWindow("VI Barcode Helper")
	controller := app.NewController(cfg)
	gui.MakeUI(window, controller, myApp)

	//window.Resize(fyne.Size{Width: 800, Height: 800})
	window.CenterOnScreen()
	window.ShowAndRun()

}
