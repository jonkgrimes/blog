package models

import (
	"github.com/gosimple/slug"
	"time"
)

type Post struct {
	Id          int64
	Title       string    `sql:"size:255"`
	Body        string    `sql:"text"`
	Slug        string    `sql:"size:128;index"`
	PublishedAt time.Time `sql:"default:NULL"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (p *Post) BeforeSave() error {
	p.CreateSlug()

	return nil
}

func (p *Post) CreateSlug() error {
	p.Slug = slug.Make(p.Title)

	return nil
}
