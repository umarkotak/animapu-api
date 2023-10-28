package sri

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/utils/render"
)

type LocationCurrentPoint struct {
	Lat       float64   `json:"latitude"`
	Long      float64   `json:"longitude"`
	Polyline  string    `json:"polyline"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *Sri) PostLocationCurrent(c *gin.Context) {
	var locationCurrentPoint LocationCurrentPoint
	c.BindJSON(&locationCurrentPoint)
	locationCurrentPoint.CreatedAt = time.Now()

	currentLogs := []LocationCurrentPoint{}
	err := s.DataGpsLog.Get(c.Request.Context(), &currentLogs)
	if err != nil {
		logrus.WithContext(c.Request.Context()).Error(err)
		render.Response(c.Request.Context(), c, map[string]string{"error": err.Error()}, nil, 422)
		return
	}

	currentLogs = append(currentLogs, locationCurrentPoint)

	err = s.DataGpsLog.Set(c.Request.Context(), currentLogs)
	if err != nil {
		logrus.WithContext(c.Request.Context()).Error(err)
		render.Response(c.Request.Context(), c, map[string]string{"error": err.Error()}, nil, 422)
		return
	}

	render.Response(c.Request.Context(), c, map[string]string{}, nil, 200)
}
