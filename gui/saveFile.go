package gui

import (
	"main/app"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type saveFile struct {
	saveFileLabel  *widget.Label
	saveFileButton *widget.Button
}

func makeSaveFile(w fyne.Window, c *app.Controller) saveFile {
	saveFileLabel := widget.NewLabel("Сохранить PDF")
	saveFileButton := widget.NewButton("Сохранить", func() {
		c.SavingFile()
	})
	return saveFile{
		saveFileLabel:  saveFileLabel,
		saveFileButton: saveFileButton,
	}
}
