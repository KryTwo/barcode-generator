package gui

import (
	"barcode-app/app"
	"barcode-app/config"
	"barcode-app/logger"
	"encoding/json"
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func makePresetSelect(w fyne.Window, c *app.Controller) fyne.Widget {
	data, err := os.ReadFile("settings.json")
	if err != nil {
		logger.LogError(err, "falied to onen settings.json")
		return nil
	}
	var settingsJSON config.JSONSettings
	json.Unmarshal(data, &settingsJSON)

	var str []string

	for _, v := range settingsJSON.Presets {
		str = append(str, v.Name)
	}

	button := widget.NewSelect(str, func(s string) {
		fmt.Printf("name select options: %v\n", s)
		c.RegeneratePreview()
		config.SetPreset(s)

	})
	button.PlaceHolder = "Выберите пресет"
	return button
}
