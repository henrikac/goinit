package main

import (
	"os"

	"github.com/henrikac/goinit/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
