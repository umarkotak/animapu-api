package config

import (
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

type (
	Config struct {
		CacheObj                    *cache.Cache
		AnimapuOnlineHost           string
		AnimapuLocalHost            string
		AnimapuGoogleServiceAccount string
		AnimapuFirebaseUrl          string
		ScrapeNinjaConfig           ScrapeNinjaConfig
		RedisConfig                 RedisConfig
	}
	ScrapeNinjaConfig struct {
		Host         string
		RapidApiHost string
		RapidApiKeys []string
	}
	RedisConfig struct {
		DbName   string
		Endpoint string
		Username string
		Password string
	}
)

var config Config

func Initialize() {
	err := godotenv.Load()
	if err != nil {
		logrus.Errorf("Error load env", err)
	}

	cacheObj := cache.New(5*time.Minute, 10*time.Minute)

	config = Config{
		CacheObj:                    cacheObj,
		AnimapuOnlineHost:           "https://animapu-api.herokuapp.com",
		AnimapuLocalHost:            "http://localhost:6001",
		AnimapuGoogleServiceAccount: os.Getenv("ANIMAPU_GOOGLE_SERVICE_ACCOUNT"),
		AnimapuFirebaseUrl:          os.Getenv("ANIMAPU_FIREBASE_URL"),

		ScrapeNinjaConfig: ScrapeNinjaConfig{
			Host:         "https://scrapeninja.p.rapidapi.com",
			RapidApiHost: "scrapeninja.p.rapidapi.com",
			RapidApiKeys: strings.Split(os.Getenv("RAPID_API_KEYS"), ","),
		},

		RedisConfig: RedisConfig{
			DbName:   "coding-free-db",
			Endpoint: "redis-19453.c292.ap-southeast-1-1.ec2.cloud.redislabs.com:19453",
			Username: "default",
			Password: os.Getenv("REDIS_PASSWORD"),
		},
	}
}

func Get() Config {
	return config
}
