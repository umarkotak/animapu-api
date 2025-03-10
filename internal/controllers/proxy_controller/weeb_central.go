package proxy_controller

import (
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/utils/render"
)

func WeebCentralImage(c *gin.Context) {
	currPath := c.Request.URL.String()
	splitPath := strings.Split(currPath, "/image_proxy/")
	if len(splitPath) != 2 {
		logrus.WithContext(c.Request.Context()).Error(models.ErrInvalidFormat)
		render.ErrorResponse(c.Request.Context(), c, models.ErrInvalidFormat, false)
		return
	}

	targetUrl := splitPath[1]

	if targetUrl == "" {
		logrus.WithContext(c.Request.Context()).Error(models.ErrInvalidTargetURL)
		render.ErrorResponse(c.Request.Context(), c, models.ErrInvalidTargetURL, false)
		return
	}

	req, err := http.NewRequest("GET", targetUrl, nil)
	if err != nil {
		logrus.WithContext(c.Request.Context()).Error(err)
		render.ErrorResponse(c.Request.Context(), c, err, false)
		return
	}

	req.Header.Set("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Set("Referer", "https://weebcentral.com/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36")
	req.Header.Set("sec-ch-ua", "\"Not(A:Brand\";v=\"99\", \"Google Chrome\";v=\"133\", \"Chromium\";v=\"133\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.WithContext(c.Request.Context()).Error(err)
		render.ErrorResponse(c.Request.Context(), c, err, false)
		return
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.WithContext(c.Request.Context()).Error(err)
		render.ErrorResponse(c.Request.Context(), c, err, false)
		return
	}

	c.Writer.Write(body)
}
