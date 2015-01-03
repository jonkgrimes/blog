package main

import (
	"log"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"gopkg.in/unrolled/render.v1"
)

func main() {
	n := negroni.New(
		negroni.NewRecovery(),
		negroni.NewLogger(),
		negroni.NewStatic(http.Dir("public")),
	)

	renderer := render.New(render.Options{
		Layout: "layout",
	})
	db := InitDb()

	c := &HomeController{Render: renderer, DB: db}
	p := &PostsController{Render: renderer, DB: db}

	router := mux.NewRouter().StrictSlash(true)
	router.Handle("/", c.Action(c.Index))
	router.Handle("/about", c.Action(c.About))

	router.PathPrefix("/posts/{id}").Subrouter()

	postRouter := router.PathPrefix("/posts/{id}").Subrouter()
	postRouter.Methods("GET").Handler(p.Action(p.Show))

	n.UseHandler(router)

	n.Run(":3000")
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
