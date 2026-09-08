package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"tooling"
)

func main() {
	var absPath string = ""
	var err error
	var qvDownloadURL string

	clear := flag.Bool("clear", false, "clear directory")

	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		files, err := os.ReadDir(".")
		tooling.Check(err)

		for _, file := range files {
			ext := strings.ToLower(filepath.Ext(file.Name()))
			if ext == ".md" || ext == ".pdf" {
				absPath, err = filepath.Abs(file.Name())
				break
			}
		}
	} else if len(args) == 1 {
		absPath, err = filepath.Abs(args[0])
		tooling.Check(err)

		ext := strings.ToLower(filepath.Ext(absPath))
		if ext != ".md" && ext != ".pdf" {
			fmt.Println("Error: Unsupported file type. Must be .md or .pdf")
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

	destination, file := Construct("quickView")

	if *clear {
		tooling.ClearPath(destination)
	}

	if !tooling.FileExists(file) {
		qvDownloadURL = GetDownloadURL()
		err = tooling.DownloadFile(destination, qvDownloadURL, "quickView.zip")
		tooling.Check(err)

		tooling.Unzip("quickView.zip", destination)
	}

	cmd := exec.Command(file, absPath)
	err = cmd.Start()
	tooling.Check(err)
}
