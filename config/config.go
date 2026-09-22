package config

import (
	"os"
	"strconv"
	"strings"
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
		RodBrowserPoolSize          int
		AdminEmails                 []string
	}
)

var (
	config Config
)

func Initialize() error {
	err := godotenv.Load()
	if err != nil {
		logrus.WithError(err).Warn("load env")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "6001"
	}

	poolSize, _ := strconv.Atoi(os.Getenv("ROD_BROWSER_POOL_SIZE"))
	if poolSize < 1 {
		poolSize = 2
	}
	adminEmailsValue := os.Getenv("ADMIN_EMAILS")
	adminEmails := strings.Split(adminEmailsValue, ",")
	if adminEmailsValue == "" {
		adminEmails = []string{"umarkotak@gmail.com"}
	}
	configuredAdminEmails := adminEmails[:0]
	for _, email := range adminEmails {
		if email = strings.TrimSpace(email); email != "" {
			configuredAdminEmails = append(configuredAdminEmails, email)
		}
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
		RodBrowserPoolSize:          poolSize,
		AdminEmails:                 configuredAdminEmails,
	}

	return nil
}

func Get() Config {
	return config
}

func IsAdminEmail(value string) bool {
	for _, email := range config.AdminEmails {
		if strings.EqualFold(value, email) {
			return true
		}
	}
	return false
}
