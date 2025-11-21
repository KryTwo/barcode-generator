package gui

import (
	"main/config"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type pSettings struct {
	label             *widget.Label
	labelMargin       *widget.Label
	labelYSpacing     *widget.Label
	labelXSpacing     *widget.Label
	labelMarginToCrop *widget.Label

	setMargin       *widget.Entry
	setYSpacing     *widget.Entry
	setXSpacing     *widget.Entry
	setMarginToCrop *widget.Entry
}

func MakePrintSettings() pSettings {
	label := widget.NewLabelWithStyle("Настройки печати", 1, fyne.TextStyle{Bold: true})

	//настройки отступов
	labelMargin := widget.NewLabel("Размер отступов (мм)")
	margin := binding.BindInt(&config.Get().Margin)
	setMargin := widget.NewEntryWithData(binding.IntToString(margin))
	setMargin.SetPlaceHolder("set margin...")

	//настройки промежутков между штрихкодами по вертикали
	labelYSpacing := widget.NewLabel("Промежуток между ШК по Y")
	ySpacing := binding.BindFloat(&config.Get().YSpacing)
	setYSpacing := widget.NewEntryWithData(binding.FloatToStringWithFormat(ySpacing, "%.0f"))
	setYSpacing.SetPlaceHolder("set spacing...")

	//настройки промежутков между штрихкодами по горизонтали
	labelXSpacing := widget.NewLabel("Промежуток между ШК по X")
	xSpacing := binding.BindFloat(&config.Get().XSpacing)
	setXSpacing := widget.NewEntryWithData(binding.FloatToStringWithFormat(xSpacing, "%.0f"))
	setXSpacing.SetPlaceHolder("set spacing...")

	//настройки линий нарезки ШК (отступы по бокам)
	labelMarginToCrop := widget.NewLabel("Отступы по бокам ШК")
	marginToCrop := binding.BindInt(&config.Get().MarginToCrop)
	setMarginToCrop := widget.NewEntryWithData(binding.IntToString(marginToCrop))
	setMarginToCrop.SetPlaceHolder("set margin to crop")

	return pSettings{
		label:             label,
		labelMargin:       labelMargin,
		labelYSpacing:     labelYSpacing,
		labelXSpacing:     labelXSpacing,
		labelMarginToCrop: labelMarginToCrop,

		setMargin:       setMargin,
		setYSpacing:     setYSpacing,
		setXSpacing:     setXSpacing,
		setMarginToCrop: setMarginToCrop,
	}
}
