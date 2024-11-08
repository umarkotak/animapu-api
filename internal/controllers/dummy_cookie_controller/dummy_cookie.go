package dummy_cookie_controller

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/utils/render"
)

type (
	DummyCookieParams struct {
		Cookies []CustomCookie `json:"cookies"`
	}

	CustomCookie struct {
		Origin string `json:"origin"` // https://chatgpt.com
		Name   string `json:"name"`   // name,
		Value  string `json:"value"`  // value,
		Path   string `json:"path"`   // "/",
		Domain string `json:"domain"` // example.com",
		Secure bool   `json:"secure"` // true,
	}
)

func HandlerDummyCookie(c *gin.Context) {
	ctx := c.Request.Context()

	var err error
	params := DummyCookieParams{}

	if c.Request.Method == "GET" {
		paramsString := c.Request.URL.Query().Get("params")
		err = json.Unmarshal([]byte(paramsString), &params)

	} else {
		err = c.BindJSON(&params)
	}
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		render.ErrorResponse(ctx, c, err, true)
		return
	}

	for _, oneCookie := range params.Cookies {
		cookie := &http.Cookie{
			Name:     oneCookie.Name,
			Value:    oneCookie.Value,
			Path:     oneCookie.Path,
			Domain:   oneCookie.Domain,
			Secure:   oneCookie.Secure,
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode,
		}
		http.SetCookie(c.Writer, cookie)
	}

	render.Response(ctx, c, map[string]any{
		"origin": params.Cookies[0].Origin,
	}, nil, 200)
}
