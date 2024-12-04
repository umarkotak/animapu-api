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

func KomikindoImage(c *gin.Context) {
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
	req.Header.Add("authority", "k7rzspb5flu6zayatfe4mh.my")
	req.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Add("accept-language", "en-US,en;q=0.9,id;q=0.8")
	req.Header.Add("cache-control", "max-age=0")
	req.Header.Add("if-modified-since", "Wed, 04 Dec 2024 09:23:30 GMT")
	req.Header.Add("if-none-match", "W/\"67501f92-52e0c\"")
	req.Header.Add("priority", "u=0, i' ")
	req.Header.Add("sec-ch-ua", "\"Google Chrome\";v=\"131\", \"Chromium\";v=\"131\", \"Not_A Brand\";v=\"24\"' ")
	req.Header.Add("sec-ch-ua-mobile", "?0' ")
	req.Header.Add("sec-ch-ua-platform", "\"macOS\"' ")
	req.Header.Add("sec-fetch-dest", "document' ")
	req.Header.Add("sec-fetch-mode", "navigate' ")
	req.Header.Add("sec-fetch-site", "none' ")
	req.Header.Add("sec-fetch-user", "?1' ")
	req.Header.Add("upgrade-insecure-requests", "1' ")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36'")

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
