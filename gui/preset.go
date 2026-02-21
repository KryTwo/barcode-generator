package gui

import (
	"barcode-app/app"
	"barcode-app/config"
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
		if c.OnSaveValidation != nil {
			c.OnSaveValidation()
		}
		if c.OnDeleteValidation != nil {
			c.OnDeleteValidation()
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

		if c.OnSaveValidation != nil {
			c.OnSaveValidation()
		}
		if c.OnDeleteValidation != nil {
			c.OnDeleteValidation()
		}
	}

	presetSelect.PlaceHolder = "Выберите пресет"

	//Проверяем на нужное состояние кнопки
	if c.OnSaveValidation != nil {
		c.OnSaveValidation()
	}
	if c.OnDeleteValidation != nil {
		c.OnDeleteValidation()
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

			if c.OnSaveValidation != nil {
				c.OnSaveValidation()
			}
			button.Refresh()
		}()
	})

	refresh := func() {
		if c.CurrentPresetName == "standart" || c.CurrentPresetName == "" {
			button.Disable()
		} else {
			button.Enable()
		}
	}

	// Устанавливаем начальное состояние
	refresh()
	c.OnSaveValidation = refresh

	return button
}

func DeletePresetButton(c *app.Controller, a fyne.App) fyne.Widget {
	var button *widget.Button

	// Состояния: 0 - обычное, 1 - ожидание подтверждения
	state := 0

	button = widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		if state == 0 {
			state = 1
			button.SetIcon(theme.WarningIcon())
			button.Disable()

			time.AfterFunc(time.Second*1, func() {
				button.Importance = widget.DangerImportance
				button.SetIcon(theme.DeleteIcon())
				button.Enable()
				button.Refresh()

				time.AfterFunc(time.Second*3, func() {
					if state == 1 {
						state = 0
						button.Importance = widget.MediumImportance
						if c.OnDeleteValidation != nil {
							c.OnDeleteValidation()
						}
						button.Refresh()
					}
				})
			})
			return
		}

		if state == 1 {
			state = 0
			c.DeletePreset(c.CurrentPresetName)

			button.Importance = widget.MediumImportance
			button.SetIcon(theme.ConfirmIcon())
			button.Disable()
			button.Refresh()

			if c.OnPresetChanged != nil {
				c.OnPresetChanged()
			}

			time.AfterFunc(time.Second*1, func() {
				button.SetIcon(theme.DeleteIcon())
				if c.OnDeleteValidation != nil {
					c.OnDeleteValidation()
				}
				button.Refresh()
			})
		}
	})

	refreshState := func() {
		if state == 0 {
			if c.CurrentPresetName == "standart" || c.CurrentPresetName == "" {
				button.Disable()
			} else {
				button.Enable()
			}
		}
	}

	c.OnDeleteValidation = refreshState
	refreshState()

	return button
}
func setupValidation(c *app.Controller, saveBtn, deleteBtn *widget.Button) {
	c.OnDeleteValidation = func() {
		isStandart := c.CurrentPresetName == "standart" || c.CurrentPresetName == ""

		if isStandart {
			saveBtn.Disable()
			deleteBtn.Disable()
		} else {
			saveBtn.Enable()
			deleteBtn.Enable()
		}

		saveBtn.Refresh()
		deleteBtn.Refresh()
	}

	c.OnDeleteValidation()
}
