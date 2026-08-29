package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"tooling"
)

func main() {
	var absPath string
	var err error

	const quickView = "quickView"
	const qvDownloadURL = "https://github.com/Implycitt/quickView/releases/tag/v1.0.0"

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
		if ext == ".md" || ext == ".pdf" {
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

	destination := tooling.ConstructPath(quickView)
	err = tooling.DownloadFile(destination, qvDownloadURL)
	tooling.Check(err)

	commandFile := tooling.Unzip("quickView.zip", destination)
	cmd := exec.Command(commandFile, absPath)

	err = cmd.Start()
	tooling.Check(err)

	fmt.Printf("Opening %s in QuickView\n", filepath.Base(absPath))
}
