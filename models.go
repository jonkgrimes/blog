package main

import (
	"fmt"
	"html/template"
	"time"

	"github.com/mholt/binding"
	"github.com/russross/blackfriday"
)

type Post struct {
	Id          int64
	Title       string `sql:"size:255"`
	Body        string `sql:"text"`
	PublishedAt time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type PostForm struct {
	Title string
	Body  string
}

type HomePageView struct {
	Posts []Post
}

type PostView struct {
	Post Post
}

func (p *Post) Slug() string {
	return fmt.Sprintf("/%d/%d/%s", p.PublishedAt.Year(), p.PublishedAt.Month(), p.Title)
}

func (p *PostView) Title() string {
	return p.Post.Title
}

func (p *PostView) PrettyCreatedAt() string {
	const layout = "Jan 2, 2006 at 3:04pm"
	return p.Post.CreatedAt.Format(layout)
}

func (p *PostView) Body() template.HTML {
	return template.HTML(blackfriday.MarkdownCommon([]byte(p.Post.Body)))
}

func (pf *PostForm) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&pf.Title: "title",
		&pf.Body:  "body",
	}
}
