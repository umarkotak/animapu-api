package affiliate_link_controller

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/repository/affiliate_link_repository"
	"github.com/umarkotak/animapu-api/internal/services/affiliate_link_service"
	"github.com/umarkotak/animapu-api/internal/utils/render"
	"github.com/umarkotak/animapu-api/internal/utils/utils"
)

type (
	AddAffiliateLinkParams struct {
		Link string `json:"link"`
	}
)

func GetRandom(c *gin.Context) {
	ctx := c.Request.Context()

	pagination := models.Pagination{
		Limit: utils.StringMustInt64(c.Query("limit")),
		Page:  utils.StringMustInt64(c.Query("page")),
	}
	pagination.SetDefault(1)

	affiliateLinks, err := affiliate_link_repository.GetRandom(ctx, pagination)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		render.ErrorResponse(ctx, c, err, true)
		return
	}

	render.Response(ctx, c, affiliateLinks, nil, 200)
}

func GetList(c *gin.Context) {
	ctx := c.Request.Context()

	pagination := models.Pagination{
		Limit: utils.StringMustInt64(c.Query("limit")),
		Page:  utils.StringMustInt64(c.Query("page")),
	}
	pagination.SetDefault(50)

	affiliateLinks, err := affiliate_link_repository.GetList(ctx, pagination)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		render.ErrorResponse(ctx, c, err, true)
		return
	}

	render.Response(ctx, c, affiliateLinks, nil, 200)
}

func AddTokopediaAffiliateLink(c *gin.Context) {
	ctx := c.Request.Context()

	params := AddAffiliateLinkParams{}

	err := c.BindJSON(&params)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		render.ErrorResponse(ctx, c, err, true)
		return
	}

	affiliateLink, err := affiliate_link_service.AddTokopediaAffiliateLink(c.Request.Context(), params.Link)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		render.ErrorResponse(ctx, c, err, true)
		return
	}

	render.Response(ctx, c, affiliateLink, nil, 200)
}
