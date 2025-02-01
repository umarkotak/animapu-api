package contract

type (
	UserMangaActivityData struct {
		Users []UserMangaActivity `json:"users"`
	}

	UserMangaActivity struct {
		VisitorID      string         `json:"visitor_id"`
		Email          string         `json:"email"`
		MangaHistories []MangaHistory `json:"manga_histories"`
	}
)
