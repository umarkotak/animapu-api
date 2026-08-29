package main

import (
	"log"

	"github.com/umarkotak/animapu-api/internal/app"
)

func main() {
	app.Initialize()
	log.Fatal(app.Start())
}
