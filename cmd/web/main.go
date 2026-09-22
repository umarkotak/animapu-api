package main

import (
	"log"
	"os"

	"github.com/umarkotak/animapu-api/config"
	"github.com/umarkotak/animapu-api/datastore"
	"github.com/umarkotak/animapu-api/internal/app"
)

func main() {
	if isMigrateUp(os.Args) {
		if err := config.Initialize(); err != nil {
			log.Fatal(err)
		}
		if err := datastore.MigrateUp(); err != nil {
			log.Fatal(err)
		}
		return
	}

	if err := app.Initialize(); err != nil {
		log.Fatal(err)
	}
	defer datastore.Close()

	if err := app.Start(); err != nil {
		log.Print(err)
	}
}

func isMigrateUp(args []string) bool {
	return len(args) == 3 && args[1] == "migrate" && args[2] == "up"
}
