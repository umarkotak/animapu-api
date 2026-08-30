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
		MangameeApiHost             string
		CollyTimeout                time.Duration
		DbUrl                       string
		RodHeadless                 bool
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
		MangameeApiHost:             os.Getenv("MANGAMEE_API_HOST"),
		CollyTimeout:                5 * time.Minute,
		DbUrl:                       os.Getenv("DB_URL"),
		RodHeadless:                 os.Getenv("ROD_HEADLESS") == "true",
	}

	return nil
}

func Get() Config {
	return config
}
