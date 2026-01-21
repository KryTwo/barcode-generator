package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func makeMenu(a fyne.App) *fyne.MainMenu {
	fileOpenFile := fyne.NewMenuItem("Open File", func() {
		fmt.Println("some")
	})
	fileOpenRecentFile := fyne.NewMenuItem("Open Recent", nil)
	fileOpenRecentFile.ChildMenu = fyne.NewMenu("", fyne.NewMenuItem("last opened", nil), fyne.NewMenuItem("first opened", nil))

	settingsItem := fyne.NewMenuItem("TODO settings", func() {
		fmt.Println("settings print")
	})

	showAbout := func() {
		w := a.NewWindow("About")
		w.CenterOnScreen()
		w.SetContent(widget.NewLabelWithStyle(aboutContent, fyne.TextAlignCenter, fyne.TextStyle{}))
		w.Show()
	}

	//проверка новой версии. Одной кнопкой? Вывести текущую версию.
	checkUpdates := func() {
		w := a.NewWindow("Version")
		w.CenterOnScreen()
		w.Resize(fyne.Size{Width: 200, Height: 100})
		content := container.NewVBox(
			widget.NewLabel("Current vesion: 0.0.1"),
			widget.NewLabel("Available version: 0.0.1"),
			widget.NewButton("Update", func() { fmt.Println("Updating...") }),
		)

		w.SetContent(container.NewCenter(content))
		w.Show()
	}
	helpAbout := fyne.NewMenuItem("About", showAbout)
	helpCheckUpdates := fyne.NewMenuItem("Check updates", checkUpdates)
	//help_help := fyne.NewMenuItem("About", showAbout)

	File := fyne.NewMenu("File", fileOpenFile, fileOpenRecentFile)
	Settings := fyne.NewMenu("Settings", settingsItem)
	help := fyne.NewMenu("Help", helpAbout, helpCheckUpdates)

	main_menu := fyne.NewMainMenu(File, Settings, help)
	return main_menu

}

var aboutContent string = "This program is designed to simplify the barcode printing process\n" +
	"and reduce paper consumption.\n\n" +
	"Specially for Vi.ru\n\n" +
	"Distributed free of charge."
