package main

import (
	"log"

	"github.com/umarkotak/animapu-api/datastore"
	"github.com/umarkotak/animapu-api/internal/app"
)

func main() {
	if err := app.Initialize(); err != nil {
		log.Fatal(err)
	}
	defer datastore.Close()

	if err := app.Start(); err != nil {
		log.Print(err)
	}
}
