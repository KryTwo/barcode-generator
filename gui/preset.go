package gui

import (
	"barcode-app/app"
	"barcode-app/config"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func makePresetSelect(c *app.Controller, bcs BCSettingsWidgetsStruct, prs PrintSettingsWidgetStruct) fyne.Widget {
	c.ReadJSON()

	var presetNames []string

	for _, v := range config.ConfigJSON.Presets {
		presetNames = append(presetNames, v.Name)
	}

	presetSelect := widget.NewSelect(presetNames, func(s string) {
		c.SetPreset(s)
		c.ReadJSON()

		c.CurrentPresetName = s
		c.RegeneratePreview()
		bcs.UpdateFields()
		prs.UpdateFields()

	})

	c.OnPresetChanged = func() {
		c.ReadJSON()

		var newNames []string
		for _, v := range config.ConfigJSON.Presets {
			newNames = append(newNames, v.Name)
		}
		presetSelect.Options = newNames
		presetSelect.SetSelected(c.CurrentPresetName)
		presetSelect.Refresh()
	}

	presetSelect.PlaceHolder = "Выберите пресет"
	return presetSelect
}

// Кнопка добавления нового пресета
// Добавить автоматическое переключение на новый пресет
func AddPresetButton(w fyne.Window, c *app.Controller, a fyne.App) fyne.Widget {
	button := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		nw := a.NewWindow("Имя пресета")

		Entry := widget.NewEntry()
		Entry.SetPlaceHolder("Введите название пресета")

		Entry.OnSubmitted = func(s string) {
			c.CreatePreset(s, nw.Close)
			if c.OnPresetChanged != nil {
				c.CurrentPresetName = s
				c.OnPresetChanged()
			}
		}

		nw.SetContent(Entry)
		nw.Resize(fyne.Size{Width: 300, Height: 40})
		nw.CenterOnScreen()
		nw.Show()
	})
	return button
}
func SavePresetButton(c *app.Controller, a fyne.App) fyne.Widget {
	button := widget.NewButtonWithIcon("", theme.DocumentSaveIcon(), func() {
		if c.CurrentPresetName == "Standart" {
			return
		}
		if c.OnPresetChanged != nil {
			c.OnPresetChanged()
		}

		c.SavePreset(c.CurrentPresetName)
	})

	//Если выбран стандартный пресет, кнопка должна быть неактивна (втч визуально)
	//Проверка на доступность записи файла
	//Бэкап настроек до успешного сохранения
	//рефреш превью
	return button
}
func DeletePresetButton(c *app.Controller, a fyne.App) fyne.Widget {
	button := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {})
	//Неактивная кнопка при активном стандартном пресете
	//Автопереключение на стандартный пресет при удалении текущего пресета
	//Двойное подтверждение удаления (цвет кнопки)
	//Обновление списка оставшихся пресетов
	//ОС после удаления
	return button
}
