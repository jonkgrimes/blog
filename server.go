package main

import (
	"log"

	"github.com/go-martini/martini"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
)

type Post struct {
	Id    int64
	Title string `sql:"size:255"`
	Body  string `sql:"text"`
}

type PostForm struct {
	Title string `form:"title" binding:"required"`
	Body  string `form:"body" binding:"required"`
}

type HomePageView struct {
	Posts []Post
}

func main() {
	m := martini.Classic()
	m.Use(render.Renderer(render.Options{
		Layout: "layout",
	}))

	db := initDb()

	m.Map(db)

	m.Get("/", getHome)

	m.Group("/posts", func(martini.Router) {
		m.Get("/new", newPost)
		m.Post("/create", binding.Bind(PostForm{}), createPost)
		m.Get("/:id", showPost)
	})

	m.Run()
}

func getHome(r render.Render, db gorm.DB) {
	var posts []Post
	db.Find(&posts)
	r.HTML(200, "home", &HomePageView{Posts: posts})
}

func showPost(params martini.Params, r render.Render, db gorm.DB) {
	post := Post{}
	db.First(&post, params["id"])
	r.HTML(200, "posts/show", post)
}

func newPost(r render.Render) {
	r.HTML(200, "posts/new", nil)
}

func createPost(postForm PostForm, r render.Render, db gorm.DB) {
	post := Post{Title: postForm.Title, Body: postForm.Body}
	db.Create(&post)
	r.Redirect("/", 301)
}

func initDb() gorm.DB {
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
