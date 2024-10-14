package models

import (
	"fmt"
	"strings"
)

type (
	QueryParams struct {
		Source            string `json:"source"`
		SourceID          string `json:"source_id"`
		SecondarySourceID string `json:"secondary_source_id"`
		Page              int64  `json:"page"`
		ChapterID         string `json:"chapter_id"`
		Title             string `json:"title"`
	}

	AnimeQueryParams struct {
		Source          string `json:"source"`            //
		SourceID        string `json:"source_id"`         //
		Page            int64  `json:"page"`              //
		WatchVersion    string `json:"watch_version"`     //
		EpisodeID       string `json:"episode_id"`        //
		Title           string `json:"title"`             //
		ReleaseYear     int64  `json:"release_year"`      //
		ReleaseSeason   string `json:"release_season"`    // fall, winter, spring, summer
		FromLocal       string `json:"from_local"`        //
		Resolution      string `json:"resolution"`        // 360p, 480p, 720p
		StreamName      string `json:"stream_name"`       //
		ManualServerOpt string `json:"manual_server_opt"` //
	}
)

func (qp *QueryParams) ToKey(page string) string {
	return fmt.Sprintf(
		"%v:%v:%v:%v:%v:%v:%v",
		qp.Source,
		qp.SourceID,
		qp.SecondarySourceID,
		qp.Page,
		qp.ChapterID,
		qp.Title,
		page,
	)
}

func (qp *QueryParams) ToFbKey(page string) string {
	tempKey := fmt.Sprintf(
		"%v:%v:%v:%v:%v:%v:%v",
		qp.Source,
		qp.SourceID,
		qp.SecondarySourceID,
		qp.Page,
		qp.ChapterID,
		qp.Title,
		page,
	)
	tempKey = strings.ReplaceAll(tempKey, ".", "dot")
	return tempKey
}

func (qp *AnimeQueryParams) ToKey(page string) string {
	return fmt.Sprintf(
		"%v:%v",
		page, qp,
	)
}
