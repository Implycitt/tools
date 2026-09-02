package template

import (
	"fmt"
	"os"
	"tooling"
)

func main() {
	fmt.Println("")
}

func copyTemplate(template string) {
	srcFile, err := os.Open(template)
	tooling.Check(err)
	defer srcFile.Close()
}
