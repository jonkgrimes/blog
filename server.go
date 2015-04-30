package main

import (
	"bytes"
	"log"
	"net/http"
	"text/template"

	"github.com/codegangsta/negroni"
	"github.com/danryan/env"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/phyber/negroni-gzip/gzip"
	"gopkg.in/unrolled/render.v1"

	"blog/controllers"
	"blog/models"
)

type Config struct {
	DbName      string `env:"key=BLOG_DATABASE_NAME default=blog_development"`
	DbUser      string `env:"key=BLOG_DATABASE_USER default=jong"`
	DbPassword  string `env:"key=BLOG_DATABASE_PASSWORD"`
	Port        string `env:"key=BLOG_PORT default=:3000"`
	Environment string `env:"key=ENVIRONMENT default=development"`
	BlogUser    string `env:"key=BLOG_USER default=fluffywolf24"`
	BlogPwd     string `env:"key=BLOG_PASSWORD default=Longhorn$2"`
}

func main() {
	config := &Config{}
	checkErr(env.Process(config), "There was a configuration error: ")

	// Setup middlewares
	n := negroni.New(
		negroni.NewRecovery(),
		negroni.NewLogger(),
		negroni.NewStatic(http.Dir("public")),
	)

	// Use Gzip in production
	if config.Environment == "production" {
		n.Use(gzip.Gzip(gzip.DefaultCompression))
	}

	renderer := render.New(render.Options{
		Layout: "layout",
	})
	db := InitDb(config)

	p := &controllers.PostsController{Render: renderer, Db: db}

	router := mux.NewRouter().StrictSlash(true)

	postRouter := router.PathPrefix("/posts").Subrouter()
	postRouter.Path("/").Methods("GET").Handler(p.Action(p.Index))
	postRouter.Path("/{id}").Methods("GET").Handler(p.Action(p.Show))

	n.UseHandler(router)

	n.Run(config.Port)
}

func InitDb(c *Config) gorm.DB {
	tmpl_str := "user={{.DbUser}}{{if .DbPassword}} password={{.DbPassword}}{{end}} dbname={{.DbName}} sslmode=disable"
	tmpl, err := template.New("connection").Parse(tmpl_str)
	var b bytes.Buffer
	err = tmpl.Execute(&b, c)
	connString := b.String()

	db, err := gorm.Open("postgres", connString)

	checkErr(err, "gorm.Open failed")

	db.AutoMigrate(&models.Post{})
	db.LogMode(true)

	return db
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
