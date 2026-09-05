package tooling

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Release struct {
	TagName string `json:"tag_name"`
}

func GetRecentTag(apiUrl string) (*Release, error) {
	req, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	Check(err)

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Do(req)
	Check(err)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("No releases found")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API Failure: %d", resp.StatusCode)
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("Failed to decode: %w", err)
	}

	return &release, nil
}
