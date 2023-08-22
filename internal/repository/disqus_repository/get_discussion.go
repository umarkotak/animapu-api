package disqus_repository

import (
	"context"
	"time"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
)

type (
	GetDiscussionParams struct{}

	GetDiscussionResponse struct{}
)

func GetDiscussion(ctx context.Context, params GetDiscussionParams) (GetDiscussionResponse, error) {
	c := colly.NewCollector()
	c.SetRequestTimeout(60 * time.Second)

	c.OnHTML("#disqus-threadData", func(e *colly.HTMLElement) {
		logrus.WithContext(ctx).Info(e.Text)
	})

	err := c.Visit("https://disqus.com/embed/comments/?base=default&f=asurascans-com-1&t_i=242992%20https%3A%2F%2Fasura.nacm.xyz%2F%3Fp%3D242992&t_u=https%3A%2F%2Fasura.nacm.xyz%2F2678590577-academys-undercover-professor-chapter-55-framework%2F&t_e=Academy%E2%80%99s%20Undercover%20Professor%20Chapter%2055%20%26%238211%3B%20Framework&t_d=Academy%E2%80%99s%20Undercover%20Professor%20Chapter%2055%20%E2%80%93%20Framework&t_t=Academy%E2%80%99s%20Undercover%20Professor%20Chapter%2055%20%26%238211%3B%20Framework&s_o=default#version=d3a7e0f9d834ec1287136e3d51e7ef82")
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return GetDiscussionResponse{}, err
	}

	return GetDiscussionResponse{}, nil
}
