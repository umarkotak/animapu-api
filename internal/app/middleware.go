package app

import (
	"context"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/repository"
	"github.com/umarkotak/animapu-api/internal/services/user_service"
	"github.com/umarkotak/animapu-api/internal/utils/common_ctx"
	"github.com/umarkotak/animapu-api/internal/utils/fiber_ctx"
	"github.com/umarkotak/animapu-api/internal/utils/render"
)

func RequestID() fiber.Handler {
	return func(c fiber.Ctx) error {
		legacy := fiber_ctx.From(c)
		reqID, _ := uuid.NewRandom()
		ctx := context.WithValue(legacy.Request.Context(), "request_id", reqID)
		legacy.Request = legacy.Request.WithContext(ctx)
		c.SetContext(ctx)
		return c.Next()
	}
}

func CORSMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		legacy := fiber_ctx.From(c)
		legacy.Header("Access-Control-Allow-Origin", "*")
		legacy.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		legacy.Header(
			"Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Animapu-User-Uid, Animapu-User-Email, X-Visitor-Id, X-From-Path",
		)
		legacy.Header("Access-Control-Allow-Credentials", "true")

		if legacy.Request.Method == "OPTIONS" {
			render.Response(legacy.Request.Context(), legacy, nil, nil, 200)
			return nil
		}
		return c.Next()
	}
}

func LogRequest() fiber.Handler {
	return func(c fiber.Ctx) error {
		legacy := fiber_ctx.From(c)
		if legacy.Request.Method != "OPTIONS" {
			logrus.Infof("%v %v", legacy.Request.Method, legacy.Request.URL.Path)
		}
		return c.Next()
	}
}

func LogVisitor() fiber.Handler {
	return func(c fiber.Ctx) error {
		legacy := fiber_ctx.From(c)
		if legacy.Request.Method != "OPTIONS" && legacy.Request.Header.Get("X-Visitor-Id") != "" {
			go repository.LogVisitor(legacy.Request.Header.Get("X-Visitor-Id"), legacy.Request.Header.Get("X-From-Path"))
		}
		return c.Next()
	}
}

func CommonCtx() fiber.Handler {
	return func(c fiber.Ctx) error {
		legacy := fiber_ctx.From(c)
		commonCtx := common_ctx.FromRequestHeader(legacy.Request)

		legacy.Set(string(common_ctx.CommonCtxKey), commonCtx)
		return c.Next()
	}
}

func RegisterUser() fiber.Handler {
	return func(c fiber.Ctx) error {
		legacy := fiber_ctx.From(c)
		commonCtxInterface, exists := legacy.Get(string(common_ctx.CommonCtxKey))

		if exists {
			commonCtx := commonCtxInterface.(common_ctx.CommonCtx)

			user, err := user_service.UpsertAndGetUser(legacy.Request.Context(), models.User{
				VisitorId: commonCtx.User.VisitorId,
				Guid:      commonCtx.User.Guid,
				Email:     commonCtx.User.Email,
			})
			if err != nil {
				render.ErrorResponse(legacy.Request.Context(), legacy, err, true)
				return nil
			}
			commonCtx.User.ID = user.ID

			legacy.Set(string(common_ctx.CommonCtxKey), commonCtx)
		}

		return c.Next()
	}
}
