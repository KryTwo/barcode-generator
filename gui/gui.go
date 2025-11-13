package gui

import (
	"fmt"
	"image"
	"io"
	"log"
	"main/app"
	"main/config"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

func MakeUI(w fyne.Window, controller *app.Controller) {
	//параметры контейнера с изображением баркода
	BCImage := canvas.NewImageFromImage(nil)
	BCSize := fyne.Size{
		Width:  float32(config.Get().Width) / 72 * 300,
		Height: float32(config.Get().Hight) / 72 * 300,
	}
	BCImage.SetMinSize(BCSize)
	BCImage.FillMode = canvas.ImageFillContain
	BCImage.ScaleMode = canvas.ImageScaleSmooth
	BCContainer := container.NewStack(BCImage)
	BCContainer.Resize(BCSize)

	//параметры контейнера с превью печати
	previewImage := canvas.NewImageFromImage(nil)
	previewImage.SetMinSize(fyne.NewSize(600, 600))
	previewImage.FillMode = canvas.ImageFillContain
	previewImage.ScaleMode = canvas.ImageScaleSmooth
	previewContainer := container.NewStack(previewImage)

	printSettings := widget.NewLabel("Настройки печати")

	//настройки ширины ШК
	labelWidth := widget.NewLabel("Ширина штрихкода")
	width := binding.BindInt(&config.Get().Width)
	setWidth := widget.NewEntryWithData(binding.IntToString(width))
	setWidth.SetPlaceHolder("set width...")

	//настройки высоты ШК
	labelHight := widget.NewLabel("Высота штрихкода")
	hight := binding.BindInt(&config.Get().Hight)
	setHight := widget.NewEntryWithData(binding.IntToString(hight))
	setHight.SetPlaceHolder("set hight...")

	//настройки размера текста
	labelFontSize := widget.NewLabel("Размер текста")
	fontSize := binding.BindInt(&config.Get().FontSize)
	setFontSize := widget.NewEntryWithData(binding.IntToString(fontSize))
	setFontSize.SetPlaceHolder("set font size...")

	//настройки штрихкода
	BCSettings := container.NewGridWithRows(6, labelWidth, setWidth, labelHight, setHight, labelFontSize, setFontSize)

	setWidth.OnSubmitted = func(text string) {
		fmt.Println(text)
		controller.SetBCWidth(text)
		previewImage.Refresh()
		BCImage.Refresh()
	}

	setHight.OnSubmitted = func(text string) {
		fmt.Println(text)
		controller.SetBCHight(text)
		previewImage.Refresh()
		BCImage.Refresh()
	}

	setFontSize.OnSubmitted = func(text string) {
		fmt.Println(text)
		controller.SetFontSize(text)
		previewImage.Refresh()
		BCImage.Refresh()
	}

	controller.OnPreviewUpdated = func(r *image.RGBA) {
		previewImage.Image = r
		BCImage.Image = *controller.CropBC(r)
		previewImage.Refresh()
		BCImage.Refresh()
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
				if result.Success == false {
					log.Fatalln("result.ProcessFile error")
				}

				controller.RegeneratePreview()

				if err != nil {
					dialog.ShowError(err, w)
					return
				}
			}, w)
			dlg.SetFilter(storage.NewExtensionFileFilter([]string{".csv"}))
			dlg.Show()
		}),
	)

	printPreview := container.NewVBox(
		widget.NewLabel("print preview"),
		previewContainer,
	)

	leftPanel := container.NewVBox(
		BCContainer,
		widget.NewSeparator(),
		BCSettings,
		widget.NewSeparator(),
		printSettings,
		widget.NewSeparator(),
		fileOpen,
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
