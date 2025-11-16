package gui

import (
	"main/config"

	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type pSettings struct {
	labelMargin   *widget.Label
	labelYSpacing *widget.Label
	labelXSpacing *widget.Label

	setMargin   *widget.Entry
	setYSpacing *widget.Entry
	setXSpacing *widget.Entry
}

func MakePrintSettings() pSettings {

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

	return pSettings{
		labelMargin:   labelMargin,
		labelYSpacing: labelYSpacing,
		labelXSpacing: labelXSpacing,

		setMargin:   setMargin,
		setYSpacing: setYSpacing,
		setXSpacing: setXSpacing,
	}
}
