package repository

import (
	"github.com/patrickmn/go-cache"
)

func GoCache() *cache.Cache {
	return goCache
}
