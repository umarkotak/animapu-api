package app

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/repository"
	"github.com/umarkotak/animapu-api/internal/services/user_service"
	"github.com/umarkotak/animapu-api/internal/utils/common_ctx"
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
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Header(
			"Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Animapu-User-Uid, Animapu-User-Email, X-Visitor-Id, X-From-Path",
		)
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			render.Response(c.Request.Context(), c, nil, nil, 200)
			return
		}
		c.Next()
	}
}

func LogRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != "OPTIONS" {
			logrus.Infof("%v %v", c.Request.Method, c.Request.URL.Path)
		}
		c.Next()
	}
}

func LogVisitor() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != "OPTIONS" && c.Request.Header.Get("X-Visitor-Id") != "" {
			go repository.LogVisitor(c)
		}
		c.Next()
	}
}

func CommonCtx() gin.HandlerFunc {
	return func(c *gin.Context) {
		commonCtx := common_ctx.FromRequestHeader(c.Request)

		c.Set(string(common_ctx.CommonCtxKey), commonCtx)
	}
}

func RegisterUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		commonCtxInterface, exists := c.Get(string(common_ctx.CommonCtxKey))

		if exists {
			commonCtx := commonCtxInterface.(common_ctx.CommonCtx)

			user, err := user_service.UpsertAndGetUser(c.Request.Context(), models.User{
				VisitorId: commonCtx.User.VisitorId,
				Guid:      commonCtx.User.Guid,
				Email:     commonCtx.User.Email,
			})
			if err != nil {
				render.ErrorResponse(c.Request.Context(), c, err, true)
				return
			}
			commonCtx.User.ID = user.ID

			c.Set(string(common_ctx.CommonCtxKey), commonCtx)
		}

		c.Next()
	}
}
