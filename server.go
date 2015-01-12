package main

import (
	"fmt"
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

	router.PathPrefix("/admin").Handler(negroni.New(
		negroni.HandlerFunc(AdminAuth),
		negroni.Wrap(adminRouter),
	))

	n.UseHandler(router)

	n.Run(":8080")
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

func AdminAuth(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	fmt.Println("Authenticate...")
	next(w, r)
}
