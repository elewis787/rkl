package main

import (
	"log"

	"github.com/elewis787/rkl/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
