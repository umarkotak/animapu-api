package proxy_controller

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/utils/render"
)

func FizmangaImage(c *gin.Context) {
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

	req.Header.Set("authority", "my4.1stkmgv2.com")
	req.Header.Set("accept", "image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8")
	req.Header.Set("accept-language", "en-US,en;q=0.9,id;q=0.8")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("cookie", "cf_clearance=sBFkc37NOXgGOf_EoykEXPvAe_4D6ell1jCgsFW65Lw-1661607502-0-250")
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("referer", "https://fizmanga.com/")
	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"104\", \" Not A;Brand\";v=\"99\", \"Google Chrome\";v=\"104\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Set("sec-fetch-dest", "image")
	req.Header.Set("sec-fetch-mode", "no-cors")
	req.Header.Set("sec-fetch-site", "cross-site")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.WithContext(c.Request.Context()).Error(err)
		render.ErrorResponse(c.Request.Context(), c, err, false)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithContext(c.Request.Context()).Error(err)
		render.ErrorResponse(c.Request.Context(), c, err, false)
		return
	}

	c.Writer.Write(body)
}
