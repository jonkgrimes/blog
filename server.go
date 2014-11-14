package main

import (
	"log"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"gopkg.in/unrolled/render.v1"
)

type key int

const DB key = 0
const Renderer key = 1

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
	n := negroni.New(
		negroni.NewRecovery(),
		negroni.HandlerFunc(SetContext),
		negroni.NewLogger(),
		negroni.NewStatic(http.Dir("public")),
	)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", HomeHandler)

	// posts collection
	// posts := router.Path("/posts").Subrouter()
	//posts.Methods("GET").HandlerFunc(PostsIndexHandler)
	//posts.Methods("POST").HandlerFunc(PostsCreateHandler)

	// posts singular
	//post := r.PathPrefix("/posts/{id}/").Subrouter()
	//post.Methods("GET").Path("/edit").HandlerFunc(PostEditHandler)
	//post.Methods("GET").HandlerFunc(PostShowHandler)
	//post.Methods("PUT", "POST").HandlerFunc(PostUpdateHandler)
	//post.Methods("DELETE").HandlerFunc(PostDeleteHandler)

	n.UseHandler(router)

	n.Run(":3000")

}

func HomeHandler(rw http.ResponseWriter, r *http.Request) {
	db := GetDB(r)
	renderer := GetRenderer(r)
	var posts []Post
	db.Find(&posts)
	renderer.HTML(rw, http.StatusOK, "home", &HomePageView{Posts: posts})
}

func SetContext(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	context.Set(r, DB, InitDb())
	context.Set(r, Renderer, render.New(render.Options{
		Layout: "layout",
	}))

	next(rw, r)
}

func GetRenderer(r *http.Request) *render.Render {
	if rv := context.Get(r, Renderer); rv != nil {
		return rv.(*render.Render)
	}
	return nil
}

func SetRenderer(r *http.Request, renderer *render.Render) {
	context.Set(r, Renderer, renderer)
}

func GetDB(r *http.Request) gorm.DB {
	if rv := context.Get(r, DB); rv != nil {
		return rv.(gorm.DB)
	}
	return gorm.DB{}
}

func SetDB(r *http.Request, db gorm.DB) {
	context.Set(r, DB, db)
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
