package gogo_anime_new

type (
	GogoAnimeServerParams struct {
		DataType          string `json:"data-type"`
		DataEncryptedUrl1 string `json:"data-encrypted-url1"`
		DataEncryptedUrl2 string `json:"data-encrypted-url2"`
		DataEncryptedUrl3 string `json:"data-encrypted-url3"`
		EpisodeUrl        string `json:"-"`
		FeatureImage      string `json:"-"`
	}

	GogoAnimeStream struct {
		Mode string `json:"mode"` // Enum: iframe, mp4
		Src  string `json:"src"`  //
	}
)
