package main

import (
	"log"

	"github.com/artemxgod/project-gen-go/internal/generator"
)

func main() {
	gener := generator.New()

	if err := gener.Generate(); err != nil {
		log.Fatal(err)
	}
}
