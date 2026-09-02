package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
	"tooling"
)

type Release struct {
	TagName string `json:"tag_name"`
}

func main() {
	var absPath string
	var err error

	var qvDownloadURL = GetDownloadURL()

	if len(os.Args) == 1 {
		files, err := os.ReadDir(".")
		tooling.Check(err)
		for _, file := range files {
			ext := strings.ToLower(filepath.Ext(file.Name()))
			if ext == ".md" || ext == ".pdf" {
				absPath, err = filepath.Abs(file.Name())
			}
		}
	} else if len(os.Args) == 2 {
		path := os.Args[1]
		absPath, err = filepath.Abs(path)

		ext := strings.ToLower(filepath.Ext(absPath))
		if ext != ".md" && ext != ".pdf" {
			fmt.Println("Error: Unsupported file type")
			os.Exit(1)
		}
	} else {
		fmt.Println("Usage: qv <file.md | file.pdf>")
		os.Exit(1)
	}
	tooling.Check(err)

	stat, err := os.Stat(absPath)
	if os.IsNotExist(err) || stat.IsDir() {
		fmt.Println("Error: File does not exist")
		os.Exit(1)
	}

	destination := tooling.ConstructPath("quickView")

	if !tooling.FileExists(destination) {
		err = tooling.DownloadFile(destination, qvDownloadURL)
		tooling.Check(err)

		tooling.Unzip("quickView.zip", destination)
	}

	var argument string
	switch opsys := runtime.GOOS; opsys {
	case "windows":
		argument = destination+"quickView.exe"
	case "linux":
		argument = destination+"quickview"
	default:
		argument = ""
	}

	cmd := exec.Command(argument, absPath)
	err = cmd.Start()
	tooling.Check(err)
}

func GetDownloadURL() (url string) {
	release, err := GetRecentTag()
	tooling.Check(err)
	tagName := release.TagName[1:]

	switch opsys := runtime.GOOS; opsys {
	case "windows":
		url = fmt.Sprintf("https://github.com/Implycitt/quickView/releases/download/v%[1]s/QuickView-%[1]s-win.zip", tagName)
	case "linux":
		url = fmt.Sprintf("https://github.com/Implycitt/quickView/releases/download/v%[1]s/QuickView-%[1]s.zip", tagName)
	case "darwin":
		url = fmt.Sprintf("https://github.com/Implycitt/quickView/releases/download/v%[1]s/QuickView-%[1]s-arm64-mac.zip", tagName)
	default:
		url = ""
	}
	return url
}

func GetRecentTag() (*Release, error) {
	var apiUrl string = "https://api.github.com/repos/Implycitt/quickView/releases/latest"
	req, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	tooling.Check(err)

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Do(req)
	tooling.Check(err)
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
