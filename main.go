package main

import (
	"os"

	"github.com/elewis787/rcl/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(0)
	}
}
