package main

import (
	"context"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/config"
)

func main() {
	ctx := context.Background()

	c := colly.NewCollector()
	c.SetRequestTimeout(config.Get().CollyTimeout)

	c.OnHTML("div", func(e *colly.HTMLElement) {
		logrus.Infof("%+v", e.Text)
	})

	err := c.Visit("https://www.mangasee123.com/")
	if err != nil {
		logrus.WithContext(ctx).Error(err)
	}
}
