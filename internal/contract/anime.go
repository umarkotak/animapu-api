package contract

import (
	"strings"
	"time"

	"github.com/umarkotak/animapu-api/internal/utils/utils"
)

type (
	Anime struct {
		ID                 string    `json:"id"`
		Source             string    `json:"source"`
		Title              string    `json:"title"`
		AltTitles          []string  `json:"alt_titles"`
		Description        string    `json:"description"`
		LatestEpisode      float64   `json:"latest_episode"`
		CoverUrls          []string  `json:"cover_urls"`
		Genres             []string  `json:"genres"`
		Episodes           []Episode `json:"episodes"`
		OriginalLink       string    `json:"original_link"`
		ReleaseMonth       string    `json:"release_month"`
		ReleaseSeason      string    `json:"release_season"`
		ReleaseSeasonIndex int64     `json:"release_season_index"`
		ReleaseYear        int64     `json:"release_year"`
		ReleaseDate        string    `json:"release_date"`
		Score              float64   `json:"score"`
		Relations          []Anime   `json:"relations"`
		Relationship       string    `json:"relationship"`
		MultipleServer     bool      `json:"multiple_server"`
		SearchTitle        string    `json:"search_title"`
	}

	AnimeDetail struct {
		AnimensionAnimeID  string          `json:"animension_anime_id"`  //
		MasterTitle        string          `json:"master_title"`         //
		MasterTitleTag     string          `json:"master_title_tag"`     // Lower case, no space, only alphabet
		AltTitle           string          `json:"alt_title"`            //
		Title              string          `json:"title"`                //
		Description        string          `json:"description"`          //
		Genres             []string        `json:"genres"`               //
		SeasonLabel        string          `json:"season_label"`         // Enum: ["", "movie"]
		SeasonIndex        int             `json:"season_index"`         // Season ordering based on title, eg: bleach s1, bleach s2, ...
		ReleaseYear        int             `json:"release_year"`         //
		ReleaseSeasonName  string          `json:"release_season_name"`  //
		ReleaseSeasonIndex int             `json:"release_season_index"` //
		MalScore           float64         `json:"mal_score"`            //
		CoverURL           string          `json:"cover_url"`            //
		CoverURLs          []string        `json:"cover_urls"`           //
		HeaderCoverURL     string          `json:"header_cover_url"`     //
		Episodes           []AnimeEpisode  `json:"episodes"`             //
		VideoSources       []VideoSource   `json:"video_sources"`        //
		TotalEpisode       int64           `json:"total_episode"`        //
		Relations          []AnimeRelation `json:"relations"`            //
		YtUrl              string          `json:"yt_url"`               //
		YtIdUrl            string          `json:"yt_id_url"`            //
		Status             string          `json:"status"`               // Enum: finished, ongoing
		Type               string          `json:"type"`                 // Enum: series, movie
		LastSyncAt         time.Time       `json:"last_sync_at"`         //
	}

	VideoSource struct {
		SourceType      string `json:"source_type"`
		SourceLabel     string `json:"source_label"`
		VideoIdentifier string `json:"video_identifier"`
	}

	AnimeEpisode struct {
		AnimensionAnimeID   string   `json:"animension_anime_id"`
		AnimensionEpisodeID string   `json:"animension_episode_id"`
		EpisodeTitle        string   `json:"episode_title"`
		EpisodeNumber       float64  `json:"episode_number"`
		CoverURL            string   `json:"cover_url"`
		CoverURLs           []string `json:"cover_urls"`
		RawHlsPlaybackURL   string   `json:"raw_hls_playback_url"`
		YtUrl               string   `json:"yt_url"`
		YtIdUrl             string   `json:"yt_id_url"`
	}

	AnimeRelation struct {
		AnimeID      string   `json:"anime_id"`
		Relationship string   `json:"relationship"` // Enum: prequel, sequel
		Title        string   `json:"title"`
		CoverUrl     string   `json:"cover_url"`
		CoverUrls    []string `json:"cover_urls"`
	}

	AnimeSummary struct {
		AnimensionAnimeID   string  `json:"animension_anime_id"`
		AnimensionEpisodeID string  `json:"animension_episode_id"`
		CoverURL            string  `json:"cover_url"`
		AvatarURL           string  `json:"avatar_url"`
		Title               string  `json:"title"`
		MasterTitle         string  `json:"master_title"`
		Status              string  `json:"status"` // Enum: completed, ongoing
		LastEpisode         float64 `json:"last_episode"`
		LastEpisodeString   string  `json:"last_episode_string"`
	}

	AnimePerSeason struct {
		ReleaseYear int64   `json:"release_year"`
		SeasonName  string  `json:"season_name"`
		SeasonIndex int64   `json:"season_index"`
		Animes      []Anime `json:"animes"`
	}
)

func (ad *AnimeDetail) GenerateDefault() {
	ad.MasterTitle = strings.ReplaceAll(strings.ToLower(ad.Title), "season", "")
	ad.MasterTitle = strings.TrimSpace(utils.RemoveNonAlphabetExceptSpace(ad.MasterTitle))

	ad.MasterTitleTag = strings.ReplaceAll(ad.MasterTitle, " ", "-")

	ad.Type = "series"
	if ad.TotalEpisode == 1 {
		ad.Type = "movie"
	}
}
