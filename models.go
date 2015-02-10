package main

import (
	"html/template"
	"regexp"
	"strings"
	"time"

	"github.com/mholt/binding"
	"github.com/russross/blackfriday"
)

type Post struct {
	Id        int64
	Title     string `sql:"size:255"`
	Body      string `sql:"text"`
	Slug      string `sql:"size:128"`
	CreatedAt time.Time
	UpdatedAt time.Time
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
