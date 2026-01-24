package updater

import (
	"barcode-app/logger"
	_ "embed"
	"encoding/json"
)

type auth struct {
	TelegramToken string `json:"telegram_token"`
	ChatID        string `json:"chat_id"`
}

//go:embed secrets.json
var authData []byte

func getAuthData() auth {
	var s auth
	err := json.Unmarshal(authData, &s)
	if err != nil {
		logger.LogError(err, "Ошибка получения аутентификационных данных")
	}
	return s
}
