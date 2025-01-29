package contract

type (
	SeasonDetail struct {
		Year        int            `json:"year"`
		SeasonName  string         `json:"season_name"`
		SeasonIndex int            `json:"season_index"`
		AnimeData   []AnimeSummary `json:"anime_data"`
	}
)
