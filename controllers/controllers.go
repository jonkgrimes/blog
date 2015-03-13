package controllers

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/mholt/binding"
	"gopkg.in/unrolled/render.v1"

	"blog/models"
)

type Action func(rw http.ResponseWriter, r *http.Request) error

type AppController struct {
}

func (c *AppController) Action(a Action) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if err := a(rw, r); err != nil {
			http.Error(rw, err.Error(), 500)
		}
	})
}

type HomeController struct {
	*render.Render
	AppController
	Db gorm.DB
}

type PostsController struct {
	*render.Render
	AppController
	Db gorm.DB
}

type AdminController struct {
	*render.Render
	AppController
	Db gorm.DB
}

func (c *HomeController) Index(rw http.ResponseWriter, r *http.Request) error {
	var posts []models.Post
	c.Db.Where("published_at IS NOT NULL").Order("published_at DESC").Find(&posts)

	c.HTML(rw, http.StatusOK, "home", &models.HomePageView{Posts: posts})
	return nil
}

func (c *HomeController) About(rw http.ResponseWriter, r *http.Request) error {
	c.HTML(rw, http.StatusOK, "about", nil)
	return nil
}

func (c *PostsController) Show(rw http.ResponseWriter, r *http.Request) error {
	post := models.Post{}

	id := mux.Vars(r)["id"]
	c.Db.First(&post, id)

	c.HTML(rw, http.StatusOK, "posts/show", &models.PostView{Post: post})

	return nil
}

func (c *AdminController) Index(rw http.ResponseWriter, r *http.Request) error {
	var posts []models.Post
	c.Db.Order("created_at DESC").Find(&posts)

	c.HTML(rw, http.StatusOK, "admin/index", &models.HomePageView{Posts: posts})
	return nil
}

func (c *AdminController) New(rw http.ResponseWriter, r *http.Request) error {
	c.HTML(rw, http.StatusOK, "admin/posts/new", nil)
	return nil
}

func (c *AdminController) Edit(rw http.ResponseWriter, r *http.Request) error {
	post := models.Post{}

	id := mux.Vars(r)["id"]
	c.Db.First(&post, id)

	c.HTML(rw, http.StatusOK, "admin/posts/edit", &post)
	return nil
}

func (c *AdminController) Create(rw http.ResponseWriter, r *http.Request) error {
	postForm := &models.PostForm{}

	errs := binding.Bind(r, postForm)
	if errs.Handle(rw) {
		return nil
	}

	post := models.Post{Title: postForm.Title, Body: postForm.Body}

	c.Db.Save(&post)

	if !c.Db.NewRecord(post) {
		http.Redirect(rw, r, "/admin/", http.StatusFound)
	} else {
		c.HTML(rw, http.StatusOK, "admin/posts/new", &post)
	}

	return nil
}

func (c *AdminController) Update(rw http.ResponseWriter, r *http.Request) error {
	post := models.Post{}
	postForm := &models.PostForm{}

	id := mux.Vars(r)["id"]
	c.Db.First(&post, id)

	err := binding.Bind(r, postForm)
	if err.Handle(rw) {
		return nil
	}

	c.Db.Model(&post).Updates(models.Post{
		Title: postForm.Title,
		Body:  postForm.Body,
	})

	http.Redirect(rw, r, "/admin", http.StatusFound)
	return nil
}

func (c *AdminController) Publish(rw http.ResponseWriter, r *http.Request) error {
	post := models.Post{}

	id := mux.Vars(r)["id"]
	c.Db.First(&post, id)

	c.Db.Model(&post).Updates(models.Post{PublishedAt: time.Now()})

	http.Redirect(rw, r, "/admin", http.StatusFound)
	return nil
}
