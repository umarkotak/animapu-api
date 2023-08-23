package manga_controller

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/repository/disqus_repository"
	"github.com/umarkotak/animapu-api/internal/utils/render"
)

func GetMangaCommentsDisqus(c *gin.Context) {
	disqusResp, err := disqus_repository.GetDiscussion(
		c.Request.Context(),
		disqus_repository.GetDiscussionParams{
			DisqusID: c.Query("disqus_id"),
		},
	)
	if err != nil {
		logrus.WithContext(c.Request.Context()).Error(err)
		render.ErrorResponse(c.Request.Context(), c, err, false)
		return
	}

	tmpObj := map[string]interface{}{}

	json.Unmarshal([]byte(disqusResp.RawJson), &tmpObj)

	render.Response(c.Request.Context(), c, tmpObj, nil, 200)
}
