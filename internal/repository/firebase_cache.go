package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/umarkotak/animapu-api/internal/models"
)

func FbSet(ctx context.Context, key string, data interface{}, expiredAt time.Time) error {
	oneCache := genericCacheRef.Child(key)
	storedData := models.FbGenericCache{
		Data:      data,
		ExpiredAt: expiredAt,
	}
	return oneCache.Set(ctx, storedData)
}

func FbGet(ctx context.Context, key string) (string, []byte, error) {
	oneCache := genericCacheRef.Child(key)
	if oneCache == nil {
		return "", []byte{}, models.ErrCacheNotFound
	}
	storedData := models.FbGenericCache{}
	oneCache.Get(ctx, &storedData)
	if storedData.ExpiredAt.Before(time.Now().UTC()) {
		return "", []byte{}, models.ErrCacheNotFound
	}
	dataByte, _ := json.Marshal(storedData.Data)
	return string(dataByte), dataByte, nil
}
