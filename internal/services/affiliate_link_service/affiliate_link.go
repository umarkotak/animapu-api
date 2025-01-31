package affiliate_link_service

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	"github.com/umarkotak/animapu-api/internal/models"
	"github.com/umarkotak/animapu-api/internal/repository/affiliate_link_repository"
)

type (
	WebMetadata struct {
		OgTitle      string
		TwitterTitle string
		OgImage      string
		TwitterImage string
		OgUrl        string

		TokopediaImage string
	}
)

func AddTokopediaAffiliateLink(ctx context.Context, link string) (models.AffiliateLink, error) {
	var err error

	affiliateLink := models.AffiliateLink{
		ShortLink: link,
	}

	_, err = affiliate_link_repository.GetByShortLink(ctx, affiliateLink.ShortLink)
	if err == nil {
		return affiliateLink, nil
	}

	affiliateLink.LongLink, err = getTokopediaLongLink(ctx, affiliateLink.ShortLink)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return models.AffiliateLink{}, err
	}

	if affiliateLink.LongLink == "" {
		err = fmt.Errorf("empty tokopedia link")
		return models.AffiliateLink{}, err
	}

	webMetadata, err := getWebMetadata(ctx, affiliateLink.LongLink)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return models.AffiliateLink{}, err
	}

	affiliateLink = models.AffiliateLink{
		ShortLink: link,
		LongLink:  affiliateLink.LongLink,
		Name:      webMetadata.OgTitle,
		ImageUrl:  webMetadata.TokopediaImage,
	}

	affiliateLink.ID, err = affiliate_link_repository.Insert(ctx, nil, affiliateLink)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return models.AffiliateLink{}, err
	}

	return affiliateLink, nil
}

func getTokopediaLongLink(ctx context.Context, shortLink string) (string, error) {
	parsedURL, err := url.Parse(shortLink)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		logrus.WithContext(ctx).Error(err)
		return "", err
	}

	resp, err := http.Get(shortLink)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logrus.WithContext(ctx).Error(err)
		return "", err
	}

	// Parse the HTML content
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return "", err
	}

	// Due to redirection, we first need to get actual tokopedia page
	longLink, _ := doc.Find("div.input-container div.sub-heading a.secondary-action").Attr("href")

	return longLink, nil
}

func getWebMetadata(ctx context.Context, link string) (WebMetadata, error) {
	webMetadata := WebMetadata{}

	client := &http.Client{}
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
	}

	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36")
	req.Header.Add("ect", "4g")
	req.Header.Add("sec-ch-ua", "\"Not A(Brand\";v=\"8\", \"Chromium\";v=\"132\", \"Google Chrome\";v=\"132\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "\"macOS\"")

	res, err := client.Do(req)
	if err != nil {
		return WebMetadata{}, err
	}
	defer res.Body.Close()

	// Parse the HTML content
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		logrus.WithContext(ctx).Error(err)
		return WebMetadata{}, err
	}

	// Extract metadata
	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		if property, exists := s.Attr("property"); exists {
			if strings.ToLower(property) == "og:image" {
				webMetadata.OgImage, _ = s.Attr("content")
				return
			}
		}
		if name, exists := s.Attr("name"); exists {
			if strings.ToLower(name) == "twitter:image:src" {
				webMetadata.TwitterImage, _ = s.Attr("content")
				return
			}
		}

		if property, exists := s.Attr("property"); exists {
			if strings.ToLower(property) == "og:title" {
				webMetadata.OgTitle, _ = s.Attr("content")
				return
			}
		}
		if name, exists := s.Attr("name"); exists {
			if strings.ToLower(name) == "twitter:title" {
				webMetadata.TwitterTitle, _ = s.Attr("content")
				return
			}
		}

		if property, exists := s.Attr("property"); exists {
			if strings.ToLower(property) == "og:url" {
				webMetadata.OgUrl, _ = s.Attr("content")
				return
			}
		}
	})

	// <link data-rh="true" rel="preload" as="image" href="https://images.tokopedia.net/img/cache/500-square/VqbcmM/2023/6/29/619f12fb-a8b1-43ff-9471-8289d55fa0ca.jpg.webp?ect=4g" crossOrigin="anonymous"/>
	doc.Find("link").Each(func(i int, s *goquery.Selection) {
		attrVal, _ := s.Attr("href")

		if strings.HasPrefix(attrVal, "https://images.tokopedia.net/img/cache") {
			webMetadata.TokopediaImage = attrVal
		}
	})

	return webMetadata, nil
}
