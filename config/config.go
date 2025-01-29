package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type (
	Config struct {
		Port                        string
		AnimapuOnlineHost           string
		AnimapuLocalHost            string
		AnimapuGoogleServiceAccount string
		AnimapuFirebaseUrl          string
		MangameeApiHost             string
		CollyTimeout                time.Duration
		DbUrl                       string
	}
)

var (
	config Config
)

func Initialize() error {
	err := godotenv.Load()
	if err != nil {
		logrus.Errorf("Error load env", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "6001"
	}

	config = Config{
		Port:                        port,
		AnimapuOnlineHost:           os.Getenv("ANIMAPU_API_HOST"),
		AnimapuLocalHost:            "http://localhost:6001",
		AnimapuGoogleServiceAccount: os.Getenv("ANIMAPU_GOOGLE_SERVICE_ACCOUNT"),
		AnimapuFirebaseUrl:          os.Getenv("ANIMAPU_FIREBASE_URL"),
		MangameeApiHost:             os.Getenv("MANGAMEE_API_HOST"),
		CollyTimeout:                5 * time.Minute,
		DbUrl:                       os.Getenv("DB_URL"),
	}

	if config.AnimapuFirebaseUrl == "" {
		panic("firebase url unset!")
	}

	return nil
}

func Get() Config {
	return config
}
