package config

import (
	"barcode-app/logger"
	"encoding/json"
	"os"
)

func readJSON() JSONSettings {
	data, err := os.ReadFile("settings.json")
	if err != nil {
		logger.LogError(err, "failed to read settings.json")
		return JSONSettings{}
	}
	var JSONSettings JSONSettings
	json.Unmarshal(data, &JSONSettings)

	return JSONSettings
}

func writeJSON() {

}
