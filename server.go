package main

import (
	"github.com/go-martini/martini"
	"github.com/jinzhu/gorm"
	"github.com/martini-contrib/render"
)

func main() {
	m := martini.Classic()
	m.Use(render.Renderer(render.Options{
		Layout: "layout",
	}))

	m.Get("/", func(r render.Render) {
		r.HTML(200, "home", nil)
	})

	m.Run()
}

func initDb() *gorm.DB {
	db, err := gorm.Open("postgres", "user=root dbname=blog_development sslmode=disable")

	return db
}
