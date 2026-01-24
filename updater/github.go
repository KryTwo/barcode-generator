package updater

import (
	"encoding/json"
	"net/http"
	"time"
)

type GitHubRelease struct {
	TagName string `json:"tag_name"`
	HTMLURL string `json:"html_url"`
}

func GetReleaseInfo() (GitHubRelease, error) {
	var gitHubInfo GitHubRelease
	client := &http.Client{Timeout: 5 * time.Second}
	req, _ := http.NewRequest("GET", "https://api.github.com/repos/krytwo/barcode-generator/releases/latest", nil)
	req.Header.Set("User-Agent", "krytwo") // Обязательно любой заголовок

	resp, err := client.Do(req)
	if err != nil {
		return gitHubInfo, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&gitHubInfo); err != nil {
		return gitHubInfo, err
	}
	return gitHubInfo, nil
}
