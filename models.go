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
