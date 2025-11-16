package gui

import (
	"main/config"

	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type BCSettingsWidgets struct {
	LabelWidth    *widget.Label //label width
	LabelHight    *widget.Label //label hight
	LabelFontSize *widget.Label //label fontSize

	SetWidth    *widget.Entry //entry width
	SetHight    *widget.Entry //entry hight
	SetFontSize *widget.Entry //entry fontSize
}

func MakeBCSettings() BCSettingsWidgets {
	labelWidth := widget.NewLabel("Ширина штрихкода (мм)")
	width := binding.BindInt(&config.Get().Width)
	setWidth := widget.NewEntryWithData(binding.IntToString(width))
	setWidth.SetPlaceHolder("set width...")

	labelHight := widget.NewLabel("Высота штрихкода (мм)")
	hight := binding.BindInt(&config.Get().Hight)
	setHight := widget.NewEntryWithData(binding.IntToString(hight))
	setHight.SetPlaceHolder("set hight...")

	labelFontSize := widget.NewLabel("Размер текста")
	fontSize := binding.BindInt(&config.Get().FontSize)
	setFontSize := widget.NewEntryWithData(binding.IntToString(fontSize))
	setFontSize.SetPlaceHolder("set font size...")

	return BCSettingsWidgets{
		LabelWidth:    labelWidth,
		LabelHight:    labelHight,
		LabelFontSize: labelFontSize,
		SetWidth:      setWidth,
		SetHight:      setHight,
		SetFontSize:   setFontSize,
	}
}
