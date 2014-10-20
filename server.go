package main

import (
	"log"

	"github.com/go-martini/martini"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/martini-contrib/render"
)

type Post struct {
	Id    int64
	Title string `sql:"size:255"`
	Body  string `sql:"text"`
}

type HomePageView struct {
	Posts []Post
}

type Context struct {
}

func main() {
	m := martini.Classic()
	m.Use(render.Renderer(render.Options{
		Layout: "layout",
	}))

	db := InitDb()

	m.Get("/", GetHome)

	m.Group("/posts", func(martini.Router) {
		m.Get("/:id", GetPost)
		m.Post("/new", NewPost)
	})

	m.Run()
}

func GetHome(r render.Render) {
	var posts []Post
	db.Find(&posts)
	r.HTML(200, "home", &HomePageView{Posts: posts})
}

func GetPost(r render.Render, params martini.Params) {
	post := &Post{}
	db.First(&post, params["id"])
	r.HTML(200, "posts/show", post)
}

func NewPost() {
	return
}

func InitDb() gorm.DB {
	db, err := gorm.Open("postgres", "user=jonkgrimes dbname=blog_development sslmode=disable")

	checkErr(err, "gorm.Open failed")

	db.AutoMigrate(&Post{})
	db.LogMode(true)

	return db
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}