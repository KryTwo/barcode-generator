package gui

import (
	"barcode-app/app"
	"barcode-app/config"
	"fmt"
	"time"

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
		if c.OnValidationUpdate != nil {
			c.OnValidationUpdate()
		}
	})

	//обновление состояния виджета
	c.OnPresetChanged = func() {
		c.ReadJSON()

		var newNames []string
		for _, v := range config.ConfigJSON.Presets {
			newNames = append(newNames, v.Name)
		}

		//актуализация списка пресетов
		presetSelect.Options = newNames

		//установка актуального пресета активным
		presetSelect.SetSelected(c.CurrentPresetName)
		presetSelect.Refresh()

		if c.OnValidationUpdate != nil {
			c.OnValidationUpdate()
		}
	}

	presetSelect.PlaceHolder = "Выберите пресет"

	//Проверяем на нужное состояние кнопки
	if c.OnValidationUpdate != nil {
		c.OnValidationUpdate()
	}

	return presetSelect
}

// Кнопка добавления нового пресета
func AddPresetButton(w fyne.Window, c *app.Controller, a fyne.App) fyne.Widget {
	button := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		nw := a.NewWindow("Имя пресета")

		Entry := widget.NewEntry()
		Entry.SetPlaceHolder("Введите название пресета")

		//При изменении создаем пресет, обновляем состояние виджета
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
	var button *widget.Button
	button = widget.NewButtonWithIcon("", theme.DocumentSaveIcon(), func() {
		c.SavePreset(c.CurrentPresetName)

		oldIcon := button.Icon
		button.SetIcon(theme.ConfirmIcon())
		button.Disable()
		button.Refresh()

		go func() {
			time.Sleep(time.Second * 2)
			button.SetIcon(oldIcon)

			if c.OnValidationUpdate != nil {
				c.OnValidationUpdate()
			}
			button.Refresh()
		}()
	})

	refresh := func() {
		if c.CurrentPresetName == "standart" || c.CurrentPresetName == "" {
			fmt.Printf("c.CurrentPresetName: %v\n", c.CurrentPresetName)
			fmt.Println("must disale")
			button.Disable()
		} else {
			fmt.Printf("c.CurrentPresetName: %v\n", c.CurrentPresetName)
			fmt.Println("must enable")
			button.Enable()
		}
	}

	// Устанавливаем начальное состояние
	refresh()
	c.OnValidationUpdate = refresh

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
