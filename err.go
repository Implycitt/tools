package tools

import (
	"fmt"
	"os"
)

func Check(err error) {
	if err != nil {
		fmt.Println("Error: %s", err.Error())
		os.Exit(1)
	}
}
