package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/mholt/binding"
	"gopkg.in/unrolled/render.v1"
)

type Action func(rw http.ResponseWriter, r *http.Request) error

type AppController struct{}

func (c *AppController) Action(a Action) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if err := a(rw, r); err != nil {
			http.Error(rw, err.Error(), 500)
		}
	})
}

type HomeController struct {
	AppController
	*render.Render
	db gorm.DB
}

type PostsController struct {
	AppController
	*render.Render
	db gorm.DB
}

type AdminController struct {
	AppController
	*render.Render
	db gorm.DB
}

func (c *HomeController) Index(rw http.ResponseWriter, r *http.Request) error {
	var posts []Post
	c.db.Order("created_at DESC").Find(&posts)

	c.HTML(rw, http.StatusOK, "home", &HomePageView{Posts: posts})
	return nil
}

func (c *HomeController) About(rw http.ResponseWriter, r *http.Request) error {
	c.HTML(rw, http.StatusOK, "about", nil)
	return nil
}

func (c *PostsController) Show(rw http.ResponseWriter, r *http.Request) error {
	post := Post{}

	id := mux.Vars(r)["id"]
	c.db.First(&post, id)

	c.HTML(rw, http.StatusOK, "posts/show", &PostView{Post: post})

	return nil
}

func (c *AdminController) Index(rw http.ResponseWriter, r *http.Request) error {
	var posts []Post
	c.db.Order("created_at DESC").Find(&posts)

	c.HTML(rw, http.StatusOK, "admin/index", &HomePageView{Posts: posts})
	return nil
}

func (c *AdminController) New(rw http.ResponseWriter, r *http.Request) error {
	c.HTML(rw, http.StatusOK, "admin/posts/new", nil)
	return nil
}

func (c *AdminController) Edit(rw http.ResponseWriter, r *http.Request) error {
	post := Post{}

	id := mux.Vars(r)["id"]
	c.db.First(&post, id)

	c.HTML(rw, http.StatusOK, "admin/posts/edit", &post)
	return nil
}

func (c *AdminController) Create(rw http.ResponseWriter, r *http.Request) error {
	postForm := &PostForm{}

	errs := binding.Bind(r, postForm)
	if errs.Handle(rw) {
		return nil
	}

	post := Post{Title: postForm.Title, Body: postForm.Body}

	c.db.Save(&post)

	if !c.db.NewRecord(post) {
		http.Redirect(rw, r, "/admin/", http.StatusFound)
	} else {
		c.HTML(rw, http.StatusOK, "admin/posts/new", &post)
	}

	return nil
}

func (c *AdminController) Update(rw http.ResponseWriter, r *http.Request) error {
	post := Post{}
	postForm := &PostForm{}

	id := mux.Vars(r)["id"]
	c.db.First(&post, id)

	err := binding.Bind(r, postForm)
	if err.Handle(rw) {
		return nil
	}

	c.db.Model(&post).Updates(Post{
		Title: postForm.Title,
		Body:  postForm.Body,
	})

	http.Redirect(rw, r, "/admin", http.StatusFound)
	return nil
}
