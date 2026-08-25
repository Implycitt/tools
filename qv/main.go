package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"tools"
)

func main() {
	if len(os.Args) < 2 {
	 	fmt.Println("Usage: qv <file.md | file.pdf>")
		os.Exit(1)
	}

	path := os.Args[1]
	absPath, err := filepath.Abs(path)
	tools.Check(err)

	stat, err := os.Stat(absPath)
	if os.IsNotExist(err) || stat.IsDir() {
		fmt.Println("Error: File does not exist")
		os.Exit(1)
	}

	ext := strings.ToLower(filepath.Ext(absPath))
	if ext != ".md" && ext != ".pdf" {
		fmt.Println("Error: Unsupported file type")
		os.Exit(1)
	}

	quickView := ""
	cmd := exec.Command(quickView, absPath)

	err = cmd.Start()
	tools.Check(err)

	fmt.Printf("Opening %s in QuickView\n", filepath.Base(absPath))
}


