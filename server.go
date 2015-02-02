package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/codegangsta/negroni"
	"github.com/danryan/env"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/nabeken/negroni-auth"
	"github.com/phyber/negroni-gzip/gzip"
	"gopkg.in/unrolled/render.v1"
)

type Config struct {
	DbName      string `env:"key=BLOG_DATABASE_NAME default=blog_development"`
	DbUser      string `env:"key=BLOG_DATABASE_USER default=jong"`
	DbPassword  string `env:"key=BLOG_DATABASE_PASSWORD"`
	Port        string `env:"key=BLOG_PORT default=:8080"`
	Environment string `env:"key=ENVIRONMENT default=development"`
	BlogUser    string `env:"key=BLOG_USER default=fluffywolf24"`
	BlogPwd     string `env:"key=BLOG_PASSWORD default=Longhorn$2"`
}

func main() {
	config := &Config{}
	if err := env.Process(config); err != nil {
		fmt.Println(err)
	}

	n := negroni.New(
		negroni.NewRecovery(),
		negroni.NewLogger(),
		negroni.NewStatic(http.Dir("public")),
	)

	if config.Environment == "production" {
		n.Use(gzip.Gzip(gzip.DefaultCompression))
	}

	renderer := render.New(render.Options{
		Layout: "layout",
	})
	db := InitDb(config)

	c := &HomeController{Render: renderer, db: db}
	p := &PostsController{Render: renderer, db: db}
	a := &AdminController{Render: renderer, db: db}

	// public routes
	router := mux.NewRouter().StrictSlash(true)
	router.Handle("/", c.Action(c.Index))
	router.Handle("/about", c.Action(c.About))

	postRouter := router.PathPrefix("/posts").Subrouter()
	postRouter.Path("/{id}").Methods("GET").Handler(p.Action(p.Show))

	// admin routes
	adminRouter := mux.NewRouter().StrictSlash(true)
	adminRouter.Path("/admin/").Handler(a.Action(a.Index))
	adminRouter.Path("/admin/posts/new").Handler(a.Action(a.New))
	adminRouter.Path("/admin/posts/{id}/edit").Handler(a.Action(a.Edit))
	adminRouter.Path("/admin/posts/{id}").Methods("POST").Handler(a.Action(a.Update))
	adminRouter.Path("/admin/posts").Methods("POST").Handler(a.Action(a.Create))

	router.PathPrefix("/admin").Handler(negroni.New(
		auth.Basic(config.BlogUser, config.BlogPwd),
		negroni.Wrap(adminRouter),
	))

	n.UseHandler(router)

	n.Run(config.Port)
}

func InitDb(c *Config) gorm.DB {
	tmpl, err := template.New("connection").Parse("user={{.DbUser}}{{if .DbPassword}} password={{.DbPassword}}{{end}} dbname={{.DbName}} sslmode=disable")
	var b bytes.Buffer
	err = tmpl.Execute(&b, c)
	connString := b.String()

	db, err := gorm.Open("postgres", connString)

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
