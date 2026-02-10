package gui

import (
	"barcode-app/app"
	"barcode-app/config"
	"barcode-app/logger"
	"encoding/json"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func makePresetSelect(c *app.Controller, bcs BCSettingsWidgets, prs PrintSettingsWidgetStruct) fyne.Widget {
	data, err := os.ReadFile("settings.json")
	if err != nil {
		logger.LogError(err, "falied to onen settings.json")
		return nil
	}

	json.Unmarshal(data, &config.ConfigJSON)

	var str []string

	for _, v := range config.ConfigJSON.Presets {
		str = append(str, v.Name)
	}

	button := widget.NewSelect(str, func(s string) {
		c.SetPreset(s)
		c.RegeneratePreview()
		bcs.UpdateFields()
		prs.UpdateFields()

	})
	button.PlaceHolder = "Выберите пресет"
	return button
}
