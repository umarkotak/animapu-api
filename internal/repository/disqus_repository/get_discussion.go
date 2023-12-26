package disqus_repository

import (
	"context"
	"fmt"
	"net/url"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/config"
)

type (
	GetDiscussionParams struct {
		DisqusID string
	}

	GetDiscussionResponse struct {
		RawJson string `json:"raw_json"`
	}
)

func GetDiscussion(ctx context.Context, params GetDiscussionParams) (GetDiscussionResponse, error) {
	if params.DisqusID == "" {
		err := fmt.Errorf("empty disqus id")
		return GetDiscussionResponse{}, err
	}

	c := colly.NewCollector()
	c.SetRequestTimeout(config.Get().CollyTimeout)

	rawJson := ""
	c.OnHTML("#disqus-threadData", func(e *colly.HTMLElement) {
		// logrus.WithContext(ctx).Info(e.Text)
		rawJson = e.Text
	})

	baseUrl, err := url.Parse("https://disqus.com/embed/comments/")
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return GetDiscussionResponse{}, err
	}

	// t_i := "242992 https://asura.nacm.xyz/?p=242992"
	// t_u := "https://asura.nacm.xyz/2678590577-academys-undercover-professor-chapter-55-framework/"
	// t_e := "Academy’s Undercover Professor Chapter 55 &#8211; Framework"
	// t_d := "Academy’s Undercover Professor Chapter 55 – Framework"
	// t_t := "Academy’s Undercover Professor Chapter 55 &#8211; Framework"

	urlParams := url.Values{}
	urlParams.Add("base", "default")
	urlParams.Add("f", "asurascans-com-1")
	urlParams.Add("t_i", QueryMustUnescape(params.DisqusID))
	// urlParams.Add("t_i", QueryMustUnescape("242992%20https%3A%2F%2Fasura.nacm.xyz%2F%3Fp%3D242992"))
	// urlParams.Add("t_u", QueryMustUnescape("https%3A%2F%2Fasura.nacm.xyz%2F2678590577-academys-undercover-professor-chapter-55-framework%2F"))
	// urlParams.Add("t_e", QueryMustUnescape("Academy%E2%80%99s%20Undercover%20Professor%20Chapter%2055%20%26%238211%3B%20Framework"))
	// urlParams.Add("t_d", QueryMustUnescape("Academy%E2%80%99s%20Undercover%20Professor%20Chapter%2055%20%E2%80%93%20Framework"))
	// urlParams.Add("t_t", QueryMustUnescape("Academy%E2%80%99s%20Undercover%20Professor%20Chapter%2055%20%26%238211%3B%20Framework"))
	urlParams.Add("s_o", "default")
	baseUrl.RawQuery = urlParams.Encode()

	// logrus.Infof("CALLING DISQUS API: %+v", baseUrl.String())
	err = c.Visit(baseUrl.String())
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return GetDiscussionResponse{}, err
	}

	return GetDiscussionResponse{
		RawJson: rawJson,
	}, nil
}

func QueryMustUnescape(u string) string {
	res, _ := url.QueryUnescape(u)
	return res
}
