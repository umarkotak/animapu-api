package models

import "fmt"

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

func (qp *QueryParams) ToKey() string {
	return fmt.Sprintf(
		"%v:%v:%v:%v:%v:%v",
		qp.Source,
		qp.SourceID,
		qp.SecondarySourceID,
		qp.Page,
		qp.ChapterID,
		qp.Title,
	)
}
