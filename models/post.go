package models

import (
	"regexp"
	"strings"
	"time"
)

type Post struct {
	Id          int64
	Title       string    `sql:"size:255"`
	Body        string    `sql:"text"`
	Slug        string    `sql:"size:128"`
	PublishedAt time.Time `sql:"default:NULL"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (p *Post) CreateSlug() string {
	entities := regexp.MustCompile("&([0-9a-z#])+;")
	and := regexp.MustCompile("&")
	dashes := regexp.MustCompile("[^a-z0-9]")
	collapse := regexp.MustCompile("-+")

	result := strings.ToLower(p.Title)
	result = strings.Trim(result, " ")
	result = entities.ReplaceAllString(result, "")
	result = and.ReplaceAllString(result, "and")
	result = dashes.ReplaceAllString(result, "-")
	result = collapse.ReplaceAllString(result, "-")

	return result
}
