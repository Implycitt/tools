package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"tooling"
)

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

// just going to have to manually update the versions if there are any future changes.
func GetDownloadURL() (url string) {
	switch opsys := runtime.GOOS; opsys {
	case "windows":
		url = "https://github.com/Implycitt/quickView/releases/download/v1.0.1/QuickView-1.0.1-win.zip"
	case "linux":
		url = "https://github.com/Implycitt/quickView/releases/download/v1.0.1/QuickView-1.0.1.zip"
	default:
		url = ""
	}
	return url
}
