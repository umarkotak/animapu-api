package contract

type (
	UserAnimeActivityData struct {
		Users []UserAnimeActivity `json:"users"`
	}

	UserAnimeActivity struct {
		VisitorID      string         `json:"visitor_id"`
		Email          string         `json:"email"`
		AnimeHistories []AnimeHistory `json:"anime_histories"`
	}
)
