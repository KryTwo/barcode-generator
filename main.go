package main

import (
	"barcode-app/app"
	"barcode-app/config"
	"barcode-app/gui"
	"barcode-app/logger"
	"os"

	fyneApp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
)

func main() {
	logger.Log.Info("start main")

	config.Init()
	cfg := config.Get()
	myApp := fyneApp.NewWithID("com.krytwo.vi-helper.bcvigen")
	myApp.Settings().SetTheme(theme.DarkTheme())
	window := myApp.NewWindow("VI Barcode Helper")
	controller := app.NewController(cfg)
	gui.MakeUI(window, controller, myApp)

	_, err := os.Stat("settings.json")
	if err != nil {
		_, err = os.Create("settings.json")

		if err != nil {
			logger.LogError(err, "settings.json creating failed")
		}
	}

	//window.Resize(fyne.Size{Width: 800, Height: 800})
	window.CenterOnScreen()
	window.ShowAndRun()

}
