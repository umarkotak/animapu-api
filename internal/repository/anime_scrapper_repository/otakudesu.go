package anime_scrapper_repository

import (
	"context"
	"fmt"
	"net/url"

	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
)

type Otakudesu struct {
	AnimapuSource   string
	Source          string
	OtakudesuHost   string
	DesusStreamHost string
}

func NewOtakudesu() Otakudesu {
	return Otakudesu{
		AnimapuSource:   models.ANIME_SOURCE_OTAKUDESU,
		Source:          "otakudesu",
		OtakudesuHost:   "https://otakudesu.wiki",
		DesusStreamHost: "https://desustream.me",
	}
}

func (s *Otakudesu) Watch(ctx context.Context, queryParams models.AnimeQueryParams) (models.EpisodeWatch, error) {
	episodeWatch := models.EpisodeWatch{}

	c := colly.NewCollector()

	iframeSrc := ""
	c.OnHTML("#pembed > div > iframe", func(e *colly.HTMLElement) {
		iframeSrc = e.Attr("src")
	})

	targetUrl := fmt.Sprintf("%v/episode/%v", s.OtakudesuHost, queryParams.EpisodeID)
	err := c.Visit(targetUrl)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return episodeWatch, err
	}
	c.Wait()

	if iframeSrc == "" {
		err = models.ErrOtakudesuFrameSourceNotFound
		logrus.WithContext(ctx).Error(err)
		return episodeWatch, err
	}

	iframeUrl, err := url.Parse(iframeSrc)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return episodeWatch, err
	}

	epId := iframeUrl.Query().Get("epId")
	if epId == "" {
		epId = iframeUrl.Query().Get("id")
	}

	if epId == "" {
		err = fmt.Errorf("frame episode id not detected")
		logrus.WithContext(ctx).Error(err)
		return episodeWatch, err
	}

	desusStreamTargetUrlHd := fmt.Sprintf(
		"%v/beta/stream/hd/?id=%v", s.DesusStreamHost, epId,
	)

	episodeWatch = models.EpisodeWatch{
		IframeUrl: desusStreamTargetUrlHd,
	}

	return episodeWatch, nil
}
