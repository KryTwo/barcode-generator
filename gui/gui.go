package gui

import (
	"fmt"
	"io"
	"log"
	"main/app"
	"main/config"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

func MakeUI(w fyne.Window, controller *app.Controller) {
	BCImage := canvas.NewImageFromImage(nil)
	BCSize := fyne.Size{
		Width:  float32(config.Get().Width) / 72 * 300,
		Height: float32(config.Get().Height) / 72 * 300,
	}
	BCImage.SetMinSize(BCSize)
	BCImage.FillMode = canvas.ImageFillContain
	BCImage.ScaleMode = canvas.ImageScaleSmooth
	BCContainer := container.NewStack(BCImage)
	BCContainer.Resize(BCSize)

	previewImage := canvas.NewImageFromImage(nil)
	previewImage.SetMinSize(fyne.NewSize(600, 600))
	previewImage.FillMode = canvas.ImageFillContain
	previewImage.ScaleMode = canvas.ImageScaleSmooth
	previewContainer := container.NewStack(previewImage)

	printSettings := widget.NewLabel("Настройки печати")
	// boolData := binding.NewBool()
	setsize := widget.NewEntry()
	setsize.SetPlaceHolder("type here...")
	BCSettings := container.NewStack(
		setsize,
	)
	setsize.OnSubmitted = func(text string) {
		fmt.Println(text)
	}
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

				BCImage.Image = *controller.CropBC(result.PreviewBC)
				BCImage.Refresh()

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
		BCContainer,
		BCSettings,
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

	w.SetContent(mainHBox)
}
