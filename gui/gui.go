package gui

import (
	"fmt"
	"image"
	"main/app"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func MakeUI(w fyne.Window, controller *app.Controller) {
	var b Barcode
	BCContainer := b.MakeBarcodePreviewContainer()
	BCSettings := MakeBCSettings()
	PrintSettings := MakePrintSettings()

	//параметры контейнера с превью печати
	previewImage := MakePrintPreview()
	previewContainer := container.NewStack(previewImage)

	//настройки штрихкода
	BCSettingsContainer := container.NewGridWithRows(6,
		BCSettings.LabelWidth,
		BCSettings.SetWidth,
		BCSettings.LabelHight,
		BCSettings.SetHight,
		BCSettings.LabelFontSize,
		BCSettings.SetFontSize,
	)

	//настройки печати
	printSettingsContainer := container.NewGridWithRows(6,
		PrintSettings.labelMargin,
		PrintSettings.setMargin,
		PrintSettings.labelXSpacing,
		PrintSettings.setXSpacing,
		PrintSettings.labelYSpacing,
		PrintSettings.setYSpacing,
	)

	setupSubmittedHandler(PrintSettings.setXSpacing, controller.SetXSpacing, previewImage, &b)
	setupSubmittedHandler(PrintSettings.setYSpacing, controller.SetYSpacing, previewImage, &b)
	setupSubmittedHandler(PrintSettings.setMargin, controller.SetMargin, previewImage, &b)
	setupSubmittedHandler(BCSettings.SetWidth, controller.SetBCWidth, previewImage, &b)
	setupSubmittedHandler(BCSettings.SetHight, controller.SetBCHight, previewImage, &b)
	setupSubmittedHandler(BCSettings.SetFontSize, controller.SetFontSize, previewImage, &b)

	controller.OnPreviewUpdated = func(r *image.RGBA) {
		previewImage.Image = r
		b.BCImage.Image = *controller.CropBC(r)
		previewImage.Refresh()
		b.BCImage.Refresh()
	}

	openFileStruct := makeOpenFile(w, controller)
	fileOpen := container.NewVBox(
		openFileStruct.openFileLabel,
		openFileStruct.openFileButton,
	)

	SaveFileContainer := makeSaveFile(w, controller)
	fileSave := container.NewVBox(
		SaveFileContainer.saveFileLabel,
		SaveFileContainer.saveFileButton,
	)

	printPreview := container.NewVBox(
		widget.NewLabel("print preview"),
		previewContainer,
	)

	leftPanel := container.NewVBox(
		BCContainer,
		widget.NewSeparator(),
		BCSettingsContainer,
		widget.NewSeparator(),
		printSettingsContainer,
		widget.NewSeparator(),
		fileOpen,
		fileSave,
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

func setupSubmittedHandler(
	entry *widget.Entry,
	handlerFunc func(string),
	previewImage *canvas.Image,
	b *Barcode,
) {
	entry.OnSubmitted = func(text string) {
		fmt.Println(text)
		handlerFunc(text)
		previewImage.Refresh()
		b.BCImage.Refresh()
	}
}
