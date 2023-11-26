package models

type (
	DbAnimension struct {
		AnimeList  []AnimeDetail           //
		AnimeMap   map[string]AnimeDetail  //
		SeasonList []SeasonDetail          //
		SeasonMap  map[string]SeasonDetail //
	}
)

func (da *DbAnimension) FindAnimeByID(id string) (AnimeDetail, bool) {
	animeDetail, found := da.AnimeMap[id]
	return animeDetail, found
}
