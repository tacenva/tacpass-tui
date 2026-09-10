package main

import (
	"log"

	"github.com/tacenva/tacpass-tui/internal/app"
)

func main() {
	if err := app.Run(false); err != nil {
		log.Fatal(err)
	}
}
