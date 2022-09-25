package models

type (
	User struct {
		Uid              string           `json:"uid"`
		ReadHistories    []Manga          `json:"read_histories"`
		ReadHistoriesMap map[string]Manga `json:"read_histories_map"`
		Libraries        []Manga          `json:"libraries"`
	}
)
