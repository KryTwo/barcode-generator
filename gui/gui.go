package gui

import (
	"fmt"
	"io"
	"log"
	"main/app"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

func MakeUI(w fyne.Window, controller *app.Controller) {
	barcodePreview := widget.NewLabel("Превью баркода")
	barcodeSettings := widget.NewLabel("Настройки баркода")
	printSettings := widget.NewLabel("Настройки печати")

	fileOpen := container.NewVBox(
		widget.NewLabel("выберите файл"),
		widget.NewButton("file", func() {
			dlg := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
				if err != nil {
					dialog.ShowError(err, w)
					return
				}
				if reader == nil {
					log.Println("cancelled")
					return
				}
				defer reader.Close()

				data, err := io.ReadAll(reader)
				if err != nil {
					dialog.ShowError(fmt.Errorf("File reading error%v\n", err), w)
					return
				}

				controller.ProcessFile(data)

				if err != nil {
					dialog.ShowError(err, w)
					return
				}
			}, w)
			dlg.SetFilter(storage.NewExtensionFileFilter([]string{".csv"}))
			dlg.Show()
		}),
	)
	printPreview := container.NewVSplit(
		widget.NewLabel("print preview"),
		widget.NewLabel("preview here"),
	)

	leftTopPanel := container.NewVSplit(
		barcodePreview,
		barcodeSettings,
	)
	leftBottomPanel := container.NewVSplit(
		printSettings,
		fileOpen,
	)
	leftPanel := container.NewVSplit(
		leftTopPanel,
		leftBottomPanel,
	)

	rightPanel := container.NewVBox(
		printPreview,
	)

	mainHBox := container.NewHSplit(
		leftPanel,
		rightPanel,
	)

	/*
		fileOpenButton := widget.NewButton("Выбрать файл", func() {
			fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
				if err != nil {
					dialog.ShowError(err, w)
					return
				}
				if reader == nil {
					log.Println("cancelled")
					return
				}
				defer reader.Close()
				fmt.Println(reader)
			}, w)
			fd.Show()
		})
	*/

	w.SetContent(mainHBox)
}
