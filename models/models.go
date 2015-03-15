package models

import (
	"html/template"

	"github.com/mholt/binding"
	"github.com/russross/blackfriday"
)

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

func (p *PostView) Title() string {
	return p.Post.Title
}

func (p *PostView) PrettyPublishedAt() string {
	const layout = "Jan 2, 2006 at 3:04pm"
	return p.Post.PublishedAt.Format(layout)
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
