package proxy_controller

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/utils/render"
)

var (
	imageCache = make(map[string][]byte)
)

func AnimensionImage(c *gin.Context) {
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

	if imgBytes, ok := imageCache[targetUrl]; ok {
		c.Writer.Write(imgBytes)
		return
	}

	req, err := http.NewRequestWithContext(c.Request.Context(), "GET", targetUrl, nil)
	if err != nil {
		logrus.WithContext(c.Request.Context()).Error(err)
		render.ErrorResponse(c.Request.Context(), c, err, false)
		return
	}
	req.Header.Add("referer", "")
	req.Header.Add("sec-ch-ua", "\"Chromium\";v=\"116\", \"Not)A;Brand\";v=\"24\", \"Google Chrome\";v=\"116\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Add("sec-fetch-dest", "image")
	req.Header.Add("sec-fetch-mode", "no-cors")
	req.Header.Add("sec-fetch-site", "cross-site")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.WithContext(c.Request.Context()).Error(err)
		render.ErrorResponse(c.Request.Context(), c, err, false)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("HTTP error %d", resp.StatusCode)
		logrus.WithContext(c.Request.Context()).Error(err)
		render.ErrorResponse(c.Request.Context(), c, err, false)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		logrus.WithContext(c.Request.Context()).Error(err)
		render.ErrorResponse(c.Request.Context(), c, err, false)
		return
	}

	imageCache[targetUrl] = body

	c.Writer.Write(body)
}
