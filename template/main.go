package template

import (
	"fmt"
	"os"
	"tooling"
)

func main() {
	fmt.println
}

func copyTemplate(template string) {
	srcFile, err := os.Open(template)
	tools.Check(err)
	defer srcFile.Close()
}
