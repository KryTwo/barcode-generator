package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func makeMenu(a fyne.App) *fyne.MainMenu {
	fileItem := fyne.NewMenuItem("TODO something", func() {
		fmt.Println("some")
	})
	settingsItem := fyne.NewMenuItem("TODO settings", func() {
		fmt.Println("settings print")
	})
	showAbout := func() {
		w := a.NewWindow("About")
		w.SetContent(widget.NewLabelWithStyle(aboutContent, fyne.TextAlignCenter, fyne.TextStyle{}))
		w.Show()
	}
	//проверка новой версии. Одной кнопкой? Вывести текущую версию.
	checkUpdates := func() {
		w := a.NewWindow("Updates")
		w.Resize(fyne.Size{Width: 200, Height: 100})
		content := container.NewVBox(
			widget.NewButton("Check", func() { fmt.Println("Last version has installed") }),
		)

		w.SetContent(container.NewCenter(content))
		w.Show()
	}
	helpAbout := fyne.NewMenuItem("About", showAbout)
	helpCheckUpdates := fyne.NewMenuItem("Check updates", checkUpdates)
	//help_help := fyne.NewMenuItem("About", showAbout)

	File := fyne.NewMenu("File", fileItem)
	Settings := fyne.NewMenu("Settings", settingsItem)
	help := fyne.NewMenu("Help", helpAbout, helpCheckUpdates)

	main_menu := fyne.NewMainMenu(File, Settings, help)
	return main_menu

}

var aboutContent string = "This program is designed to simplify the barcode printing process\nand reduce paper consumption.\n\nSpecially for Vi.ru\n\nDistributed free of charge."
