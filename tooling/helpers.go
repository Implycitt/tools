package tooling

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Check(err error) {
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		os.Exit(1)
	}
}

func Panic(err error) {
	if err != nil {
		panic(err)
	}
}

func RunCommand(command string, args []string) ([]string, error) {
	cmd := exec.Command(command, args...)

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(out), "\n")
	return lines, nil
}
