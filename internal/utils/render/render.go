package render

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
)

func Response(ctx context.Context, c *gin.Context, bodyPayload interface{}, err interface{}, status int) {
	success := true
	if status != 200 {
		success = false
	}

	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Header("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.JSON(status, gin.H{
		"success": success,
		"data":    bodyPayload,
		"error":   err,
	})
}

func ErrorResponse(ctx context.Context, c *gin.Context, err error) {
	animapuError, ok := models.ERROR_MAP[err]
	if !ok {
		err = models.ErrInternal
		animapuError = models.ERROR_MAP[err]
	}

	logrus.WithContext(ctx).Error(err)
	Response(ctx, c, map[string]string{}, animapuError, animapuError.StatusCode)
}
