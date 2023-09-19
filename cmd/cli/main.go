package main

import (
	"context"
	"time"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()

	c := colly.NewCollector()
	c.SetRequestTimeout(60 * time.Second)

	c.OnHTML("div", func(e *colly.HTMLElement) {
		logrus.Infof("%+v", e.Text)
	})

	err := c.Visit("https://www.mangasee123.com/")
	if err != nil {
		logrus.WithContext(ctx).Error(err)
	}
}
