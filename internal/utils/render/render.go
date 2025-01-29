package render

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
)

func Response(ctx context.Context, c *gin.Context, bodyPayload any, err any, status int) {
	success := true
	if status != 200 {
		success = false
	}

	// logrus.Infof("BODY RESPONSE: %+v", bodyPayload)

	if c.Request.URL.Path == "/dummy-cookie" {
		c.Header("Access-Control-Allow-Origin", c.Request.URL.Query().Get("origin"))
	} else {
		c.Header("Access-Control-Allow-Origin", "*")
	}

	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Header(
		"Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Animapu-User-Uid, Animapu-User-Email, X-Visitor-Id, X-From-Path",
	)
	c.Header("Access-Control-Allow-Credentials", "true")
	c.JSON(status, gin.H{
		"success": success,
		"data":    bodyPayload,
		"error":   err,
	})
}

func ErrorResponse(ctx context.Context, c *gin.Context, err error, showErr bool) {
	animapuError, ok := models.ERROR_MAP[err]
	if !ok {
		err = models.ErrInternal
		animapuError = models.ERROR_MAP[err]
	}

	if showErr {
		animapuError.RawError = err.Error()
	}

	logrus.WithContext(ctx).Error(err)
	Response(ctx, c, map[string]string{}, animapuError, animapuError.StatusCode)
}
