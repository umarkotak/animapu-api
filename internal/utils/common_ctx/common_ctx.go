package common_ctx

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/umarkotak/animapu-api/internal/models"
)

type (
	CommonCtxKeyType string

	CommonCtx struct {
		User models.User
	}
)

var (
	CommonCtxKey = CommonCtxKeyType("common_ctx")
)

func GetFromGinCtx(c *gin.Context) CommonCtx {
	commonCtx := CommonCtx{
		User: models.User{},
	}

	v, _ := c.Get(string(CommonCtxKey))

	if v == nil {
		return commonCtx
	}

	commonCtx, ok := v.(CommonCtx)

	if !ok {
		return commonCtx
	}

	return commonCtx
}

func FromRequestHeader(r *http.Request) CommonCtx {
	return CommonCtx{
		User: models.User{
			VisitorId: r.Header.Get("X-Visitor-Id"),
			Guid:      sql.NullString{r.Header.Get("Animapu-User-Uid"), true},
			Email:     sql.NullString{r.Header.Get("Animapu-User-Email"), true},
		},
	}
}
