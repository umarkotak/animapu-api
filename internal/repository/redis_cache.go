package repository

import (
	"github.com/go-redis/redis/v8"
)

func Redis() *redis.Client {
	return rdb
}
