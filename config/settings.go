package config

import (
	"barcode-app/logger"
	"barcode-app/structs"
	"encoding/json"
	"os"
)

type JSONSettings struct {
	LastOpenedFiles []string
	DefaultSettings *structs.Config
	Presets         []Preset
}

type Preset struct {
	Name    string
	Setting structs.Config
}

var ConfigJSON JSONSettings

func configPresetSave(name string) {
	data, err := os.ReadFile("settings.json")
	if err != nil {
		logger.LogError(err, "failed to read settings.json")
		return
	}

	json.Unmarshal(data, &ConfigJSON)

	var prs Preset
	prs.Name = name
	prs.Setting = *instance

	ConfigJSON.Presets = append(ConfigJSON.Presets, prs)

}

// первая инициализация базовыми настройками программы (только при отсутствии файла)
func configInitJson(cfg *structs.Config) {
	JSONSettings := JSONSettings{
		LastOpenedFiles: []string{},
		DefaultSettings: cfg,
		Presets: []Preset{
			{Name: "Стандартный", Setting: *cfg},
		},
	}
	//если файл уже создан, не трограем его
	_, err := os.Stat("settings.json")
	if err == nil {
		return
	}

	file, err := os.Create("settings.json")

	if err != nil {
		logger.LogError(err, "settings.json creating failed")
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(JSONSettings); err != nil {
		logger.LogError(err, "fail to encode JSONSettings to settings.json")
	}
}
