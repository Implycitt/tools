package main

import (
	"flag"
	"fmt"
)

func main() {
	Parse()
}

func Parse() {
	initPtr := flag.Bool("init", false, "initialize")
	projectPtr := flag.String("", "", "initialize")

	flag.Parse()

	if *initPtr {
		if *projectPtr == "" {
			fmt.Println("This works fine")
		}
	}
}
