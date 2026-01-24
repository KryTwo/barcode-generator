package updater

import (
	"barcode-app/logger"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func SendFeedback(message string) error {

	authData := getAuthData()

	// Формируем текст сообщения
	fullText := fmt.Sprintf("📩 НОВЫЙ ОТЗЫВ\n\n<b>От:</b> \n<b>Сообщение:</b>\n%s", message)

	// Кодируем параметры для URL
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s", authData.TelegramToken+"/sendMessage")

	data := url.Values{}
	data.Set("chat_id", authData.ChatID)
	data.Set("text", fullText)
	data.Set("parse_mode", "HTML") // Чтобы работал жирный шрифт

	// Отправляем запрос
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.PostForm(apiURL, data)
	if err != nil {
		logger.LogError(err, "не могу отправить запрос")
		return err
	}
	// body, _ := io.ReadAll(resp.Body)
	// fmt.Println("Ответ Telegram:", string(body))
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.LogError(errors.New(strconv.Itoa(resp.StatusCode)), "Сервер возвратил ошибку")
		return err
	}

	return err
}
