package gui

import (
	"barcode-app/app"
	"barcode-app/config"
	"barcode-app/logger"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type BCSettingsWidgetsStruct struct {
	Label         *widget.Label
	LabelWidth    *widget.Label //label width
	LabelHight    *widget.Label //label hight
	LabelFontSize *widget.Label //label fontSize

	SetWidth    *widget.Entry //entry width
	SetHight    *widget.Entry //entry hight
	SetFontSize *widget.Entry //entry fontSize

	SetTextWrapping *widget.Check //check textWrapping

	WidthBinding    binding.Int
	HeightBinding   binding.Int
	FontSizeBinding binding.Int
}

func MakeBCSettings(c *app.Controller) BCSettingsWidgetsStruct {
	label := widget.NewLabelWithStyle("Настройки ШК", 1, fyne.TextStyle{Bold: true})

	labelWidth := widget.NewLabel("Ширина штрихкода (мм)")
	width := binding.BindInt(&config.Get().Width)
	setWidth := widget.NewEntryWithData(binding.IntToString(width))
	setWidth.SetPlaceHolder("set width...")

	labelHight := widget.NewLabel("Высота штрихкода (мм)")
	height := binding.BindInt(&config.Get().Height)
	setHight := widget.NewEntryWithData(binding.IntToString(height))
	setHight.SetPlaceHolder("set hight...")

	labelFontSize := widget.NewLabel("Размер текста")
	fontSize := binding.BindInt(&config.Get().FontSize)
	setFontSize := widget.NewEntryWithData(binding.IntToString(fontSize))
	setFontSize.SetPlaceHolder("set font size...")

	boolData := binding.NewBool()
	boolData.Set(true)
	setTextWrapping := widget.NewCheckWithData("Перенос текста", boolData)
	listener := binding.NewDataListener(
		func() {
			checked, err := boolData.Get()
			if err != nil {
				logger.LogError(err, "failed get boolData in bcSettings")
			}
			c.SetTextWrapping(checked)
		},
	)
	boolData.AddListener(listener)

	return BCSettingsWidgetsStruct{
		Label:           label,
		LabelWidth:      labelWidth,
		LabelHight:      labelHight,
		LabelFontSize:   labelFontSize,
		SetWidth:        setWidth,
		SetHight:        setHight,
		SetFontSize:     setFontSize,
		SetTextWrapping: setTextWrapping,

		WidthBinding:    width,
		HeightBinding:   height,
		FontSizeBinding: fontSize,
	}
}
func (s *BCSettingsWidgetsStruct) UpdateFields() {
	conf := config.Get()

	s.WidthBinding.Set(conf.Width)
	s.HeightBinding.Set(conf.Height)
	s.FontSizeBinding.Set(conf.FontSize)
}
