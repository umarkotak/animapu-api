package config

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type Config struct {
	CacheObj          *cache.Cache
	AnimapuOnlineHost string
	AnimapuLocalHost  string
}

var config Config

func Initialize() {
	cacheObj := cache.New(5*time.Minute, 10*time.Minute)

	config = Config{
		CacheObj:          cacheObj,
		AnimapuOnlineHost: "https://animapu-api.herokuapp.com",
		AnimapuLocalHost:  "http://localhost:6001",
	}
}

func Get() Config {
	return config
}
