package proxy_controller

import (
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/utils/render"
)

func GenericImage(c *gin.Context) {
	referer := c.Request.URL.Query().Get("referer")
	targetUrl := c.Request.URL.Query().Get("target")

	if targetUrl == "" || referer == "" {
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
	req.Header.Set("Sec-Ch-Ua", "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"98\", \"Google Chrome\";v=\"98\"")
	req.Header.Set("Referer", referer)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")

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
