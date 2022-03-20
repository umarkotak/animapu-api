package app

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/umarkotak/animapu-api/internal/utils/render"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID, _ := uuid.NewRandom()
		ctx := context.WithValue(c.Request.Context(), "request_id", reqID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			render.Response(c.Request.Context(), c, nil, nil, 200)
			return
		}
		c.Next()
	}
}
