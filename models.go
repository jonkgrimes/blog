package main

import (
	"time"

	"github.com/russross/blackfriday"
)

type Post struct {
	Id        int64
	Title     string `sql:"size:255"`
	Body      string `sql:"text"`
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

func (p *PostView) Title() string {
	return p.Post.Title
}

func (p *PostView) PrettyCreatedAt() string {
	const layout = "Jan 2, 2006 at 3:04pm"
	return p.Post.CreatedAt.Format(layout)
}

func (p *PostView) Body() string {
	return string(blackfriday.MarkdownCommon([]byte(p.Post.Body)))
}
