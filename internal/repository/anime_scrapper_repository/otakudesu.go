package anime_scrapper_repository

type Otakudesu struct {
	Host    string
	Source  string
	ImgHost string
}

func NewOtakudesu() Otakudesu {
	return Otakudesu{
		Host:    "https://www.mangasee123.com",
		Source:  "mangasee",
		ImgHost: "https://temp.compsci88.com",
	}
}
