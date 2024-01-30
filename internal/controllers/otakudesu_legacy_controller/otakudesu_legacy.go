package otakudesu_legacy_controller

import (
	"github.com/gin-gonic/gin"
	"github.com/umarkotak/animapu-api/internal/services/anime_scrapper_service"
)

// ScrapOtakudesuAllAnimes

func HandlerAnimensionQuickScrap(c *gin.Context) {
	anime_scrapper_service.ScrapOtakudesuAllAnimes(c)

	c.JSON(200, gin.H{"message": "ok", "data": "ok"})
}
