package gui

import (
	"barcode-app/app"
	"barcode-app/layout"
	"fmt"
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func MakeUI(w fyne.Window, controller *app.Controller, a fyne.App) {
	var b Barcode
	BCContainer := b.MakeBarcodePreviewContainer()
	BCSettings := MakeBCSettings(controller)
	PrintSettings := MakePrintSettings()
	//Меню бар
	setupMenu(a, w, controller)
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

	setupSubmittedHandler(PrintSettings.setXSpacing, controller.SetXSpacing, previewImage, &b, layout.ValidateXSpacing)
	setupSubmittedHandler(PrintSettings.setYSpacing, controller.SetYSpacing, previewImage, &b, layout.ValidateYSpacing)
	setupSubmittedHandler(PrintSettings.setMargin, controller.SetMargin, previewImage, &b, layout.ValidateMargin)
	setupSubmittedHandler(PrintSettings.setMarginToCrop, controller.SetMarginToCrop, previewImage, &b, layout.ValidateMarginToCrop)
	setupSubmittedHandler(BCSettings.SetWidth, controller.SetBCWidth, previewImage, &b, layout.ValidateBCWidth)
	setupSubmittedHandler(BCSettings.SetHight, controller.SetBCHeight, previewImage, &b, layout.ValidateBCHight)
	setupSubmittedHandler(BCSettings.SetFontSize, controller.SetFontSize, previewImage, &b, layout.ValidateFontSize)

	controller.OnPreviewUpdated = func(r *image.RGBA) {
		previewImage.Image = r
		b.BCImage.Image = *controller.CropBC(r)
		previewImage.Refresh()
		b.BCImage.Refresh()
	}

	//Кнопка Выбор файла
	openFileLabel := widget.NewLabel("Выберите файл (.csv)")
	openFileButton := widget.NewButton("Файл", makeOpenFile(w, controller))
	fileOpen := container.NewVBox(
		openFileLabel,
		openFileButton,
	)

	//Кнопка сохранить
	SaveFileContainer := makeSaveFile(w, controller)
	fileSave := container.NewVBox(
		SaveFileContainer.saveFileLabel,
		SaveFileContainer.saveFileButton,
	)

	presetSelect := makePresetSelect(controller, BCSettings, PrintSettings)

	//Превью печати
	printPreview := container.NewVBox(
		widget.NewLabelWithStyle("Предпросмотр печати", 1, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		previewContainer,
	)
	//Левая панель
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
		widget.NewSeparator(),
		container.NewHBox(
			presetSelect,
			AddPresetButton(w, controller, a),
			SavePresetButton(controller, a),
			DeletePresetButton(controller, a),
		),
	)
	//Правая панель
	rightPanel := container.NewVBox(
		printPreview,
	)

	//Главное окно
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
	validateFunc func(string) bool,
) {
	entry.OnSubmitted = func(text string) {
		fmt.Println(text)

		//ограничение вводимых символов цифрами
		for _, r := range text {
			if r < '0' || r > '9' {
				return
			}
		}
		//если поле ввода пустое
		if text == "" {
			return
		}
		//если не проходим валидность по лимитам
		if !validateFunc(text) {
			return
		}

		handlerFunc(text)
		previewImage.Refresh()
		b.BCImage.Refresh()
	}
}
