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
	BCSettings := MakeBCSettings(controller)
	PrintSettings := MakePrintSettings()

	//параметры контейнера с превью печати
	previewImage := MakePrintPreview()
	previewContainer := container.NewStack(previewImage)

	//настройки штрихкода
	BCSettingsContainer := container.NewGridWithRows(9,
		BCSettings.Label,
		BCSettings.LabelWidth,
		BCSettings.SetWidth,
		BCSettings.LabelHight,
		BCSettings.SetHight,
		BCSettings.LabelFontSize,
		BCSettings.SetFontSize,
		BCSettings.SetTextWrapping,
	)

	//настройки печати
	printSettingsContainer := container.NewGridWithRows(9,
		PrintSettings.label,
		PrintSettings.labelMargin,
		PrintSettings.setMargin,
		PrintSettings.labelXSpacing,
		PrintSettings.setXSpacing,
		PrintSettings.labelYSpacing,
		PrintSettings.setYSpacing,
		PrintSettings.labelMarginToCrop,
		PrintSettings.setMarginToCrop,
	)

	setupSubmittedHandler(PrintSettings.setXSpacing, controller.SetXSpacing, previewImage, &b)
	setupSubmittedHandler(PrintSettings.setYSpacing, controller.SetYSpacing, previewImage, &b)
	setupSubmittedHandler(PrintSettings.setMargin, controller.SetMargin, previewImage, &b)
	setupSubmittedHandler(PrintSettings.setMarginToCrop, controller.SetMarginToCrop, previewImage, &b)
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
		widget.NewLabelWithStyle("Предпросмотр печати", 1, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		previewContainer,
	)

	leftPanel := container.NewVBox(
		container.NewCenter(BCContainer),
		widget.NewSeparator(),

		container.NewHBox(
			BCSettingsContainer,
			widget.NewSeparator(),
			printSettingsContainer,
		),
		widget.NewSeparator(),
		container.NewHBox(
			fileOpen,
			fileSave,
		),
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
