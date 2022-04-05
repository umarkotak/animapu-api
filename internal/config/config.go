package config

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type Config struct {
	CacheObj *cache.Cache
}

var config Config

func Initialize() {
	cacheObj := cache.New(5*time.Minute, 10*time.Minute)

	config = Config{
		CacheObj: cacheObj,
	}
}

func Get() Config {
	return config
}
