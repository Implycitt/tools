package main

import (
	"os"
	"tools"
)

func main() {
}

func copyTemplate(template string) {
	srcFile, err := os.Open(template)
	tools.Check(err)
	defer srcFile.Close()
}
