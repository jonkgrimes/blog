package main

import "time"

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

func (p Post) PrettyCreatedAt() string {
	const layout = "Jan 2, 2006 at 3:04pm"
	return p.CreatedAt.Format(layout)
}
