package models

import (
	"github.com/gosimple/slug"
	"time"
)

type Post struct {
	Id          int64     `json:"id"`
	Title       string    `sql:"size:255" json:"title"`
	Body        string    `sql:"text" json:"body"`
	Slug        string    `sql:"size:128;index" json:"slug"`
	PublishedAt time.Time `sql:"default:NULL" json:"published_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (p *Post) BeforeSave() error {
	p.CreateSlug()

	return nil
}

func (p *Post) CreateSlug() error {
	p.Slug = slug.Make(p.Title)

	return nil
}
