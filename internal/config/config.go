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
		CacheObj                        *cache.Cache
		AnimapuOnlineHost               string
		AnimapuLocalHost                string
		AnimapuGoogleServiceAccount     string
		AnimapuFirebaseUrl              string
		ScrapeNinjaConfig               ScrapeNinjaConfig
		RedisConfig                     RedisConfig
		MangameeApiHost                 string
		SriFirebaseGoogleServiceAccount string
		SriFirebaseUrl                  string
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

func Initialize() error {
	err := godotenv.Load()
	if err != nil {
		logrus.Errorf("Error load env", err)
	}

	cacheObj := cache.New(5*time.Minute, 10*time.Minute)

	animapuOnlineHost := "https://animapu-api.herokuapp.com"
	if os.Getenv("ANIMAPU_ONLINE_HOST") != "" {
		animapuOnlineHost = os.Getenv("ANIMAPU_ONLINE_HOST")
	}

	config = Config{
		AnimapuOnlineHost:           animapuOnlineHost,
		CacheObj:                    cacheObj,
		AnimapuLocalHost:            "http://localhost:6001",
		AnimapuGoogleServiceAccount: os.Getenv("ANIMAPU_GOOGLE_SERVICE_ACCOUNT"),
		AnimapuFirebaseUrl:          os.Getenv("ANIMAPU_FIREBASE_URL"),
		MangameeApiHost:             os.Getenv("MANGAMEE_API_HOST"),

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

		SriFirebaseGoogleServiceAccount: os.Getenv("SRI_FIREBASE_GOOGLE_SERVICE_ACCOUNT"),
		SriFirebaseUrl:                  os.Getenv("SRI_FIREBASE_URL"),
	}

	if config.AnimapuFirebaseUrl == "" {
		panic("firebase url unset!")
	}

	return nil
}

func Get() Config {
	return config
}
