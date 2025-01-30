package contract

import (
	"fmt"
	"strings"
)

type (
	Manga struct {
		ID                  string       `json:"id"`
		SourceID            string       `json:"source_id"`
		Source              string       `json:"source"`
		Title               string       `json:"title"`
		Description         string       `json:"description"`
		Genres              []string     `json:"genres"`
		Status              string       `json:"status"`
		Rating              string       `json:"rating"`
		LatestChapterID     string       `json:"latest_chapter_id"`
		LatestChapterNumber float64      `json:"latest_chapter_number"`
		LatestChapterTitle  string       `json:"latest_chapter_title"`
		ChapterPaginated    bool         `json:"chapter_paginated"`
		Chapters            []Chapter    `json:"chapters"`
		CoverImages         []CoverImage `json:"cover_image"`
		PopularityPoint     int64        `json:"popularity_point"`
		ReadCount           int64        `json:"read_count"`
		Star                bool         `json:"star"`
		LastChapterRead     float64      `json:"last_chapter_read"`
		LastLink            string       `json:"last_link,omitempty"`
		IsInLibrary         bool         `json:"is_in_library"`
		Weight              int64        `json:"weight"`
		FollowCount         int64        `json:"follow_count"`
	}

	MangaHistory struct {
		Source              string       `json:"source"`
		SourceID            string       `json:"source_id"`
		Title               string       `json:"title"`
		LatestChapterNumber float64      `json:"latest_chapter_number"`
		CoverImages         []CoverImage `json:"cover_image"`
		LastChapterRead     float64      `json:"last_chapter_read"`
		LastLink            string       `json:"last_link"`
		IsInLibrary         bool         `json:"is_in_library"`
	}

	MangaLibrary struct {
		Source              string       `json:"source"`
		SourceID            string       `json:"source_id"`
		Title               string       `json:"title"`
		LatestChapterNumber float64      `json:"latest_chapter_number"`
		CoverImages         []CoverImage `json:"cover_image"`
		LastChapterRead     float64      `json:"last_chapter_read"`
		LastLink            string       `json:"last_link"`
		IsInLibrary         bool         `json:"is_in_library"`
	}

	CoverImage struct {
		Index     int64    `json:"index"`
		ImageUrls []string `json:"image_urls"`
	}
)

func (m *Manga) GetUniqueKey() string {
	return fmt.Sprintf(
		"%v:%v", m.Source, m.SourceID,
	)
}

func (m *Manga) GetFbUniqueKey() string {
	return strings.ReplaceAll(fmt.Sprintf(
		"%v:%v", m.Source, m.SourceID,
	), ".", "dot")
}

func (m *Manga) GenerateLatestChapter() {
	if len(m.Chapters) > 0 {
		m.LatestChapterID = m.Chapters[0].ID
		m.LatestChapterNumber = m.Chapters[0].Number
		m.LatestChapterTitle = m.Chapters[0].Title
	}
}

func (m *Manga) ImageURLs() []string {
	imageURLs := []string{}

	for _, coverImage := range m.CoverImages {
		imageURLs = append(imageURLs, coverImage.ImageUrls...)
	}

	return imageURLs
}
