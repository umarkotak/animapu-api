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
