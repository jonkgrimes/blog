package main

import (
	"log"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/mholt/binding"
	"gopkg.in/unrolled/render.v1"
)

type key int

const DB key = 0
const Renderer key = 1

type HomeController struct {
	AppController
	*render.Render
}

func main() {
	n := negroni.New(
		negroni.NewRecovery(),
		negroni.HandlerFunc(SetContext),
		negroni.NewLogger(),
		negroni.NewStatic(http.Dir("public")),
	)

	c := &HomeController{Render: render.New(render.Options{
		Layout: "layout",
	})}

	router := mux.NewRouter().StrictSlash(true)
	router.Handle("/", c.Action(c.Index))
	router.Handle("/about", c.Action(c.About))

	n.UseHandler(router)

	n.Run(":3000")
}

func (c *HomeController) Index(rw http.ResponseWriter, r *http.Request) error {
	var posts []Post
	// c.Find(&posts)

	c.HTML(rw, http.StatusOK, "home", &HomePageView{Posts: posts})
	return nil
}

func (c *HomeController) About(rw http.ResponseWriter, r *http.Request) error {
	c.HTML(rw, http.StatusOK, "about", nil)
	return nil
}

func NewPostHandler(rw http.ResponseWriter, r *http.Request) {
	renderer := GetRenderer(r)

	renderer.HTML(rw, http.StatusOK, "posts/new", nil)
}

func CreatePostHandler(rw http.ResponseWriter, r *http.Request) {
	db := GetDB(r)
	postForm := new(PostForm)
	errs := binding.Bind(r, postForm)
	if errs.Handle(rw) {
		return
	}

	post := Post{Title: postForm.Title, Body: postForm.Body}

	db.Create(&post)

	http.Redirect(rw, r, "/", 301)
}

func (pf *PostForm) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&pf.Title: binding.Field{
			Form:     "title",
			Required: true,
		},
		&pf.Body: "body",
	}
}

func ShowPostHandler(rw http.ResponseWriter, r *http.Request) {
	db := GetDB(r)
	renderer := GetRenderer(r)
	post := Post{}

	id := mux.Vars(r)["id"]
	db.First(&post, id)

	renderer.HTML(rw, http.StatusOK, "posts/show", &post)
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
