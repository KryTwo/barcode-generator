package gui

import (
	"fmt"
	"io"
	"log"
	"main/app"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

func MakeUI(w fyne.Window, controller *app.Controller) {
	previewImage := canvas.NewImageFromImage(nil)
	previewImage.SetMinSize(fyne.NewSize(600, 600))
	previewImage.FillMode = canvas.ImageFillContain
	previewImage.ScaleMode = canvas.ImageScaleSmooth
	previewContainer := container.NewStack(previewImage)

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

				result := controller.ProcessFile(data)
				if result.PreviewPNG == nil {
					log.Fatalln("result.PreviewPNG == nil")
				}

				previewImage.Image = result.PreviewPNG
				previewImage.Refresh()

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
		previewContainer,
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
