package controllers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"gopkg.in/unrolled/render.v1"

	"blog/models"
)

type PostsController struct {
	*render.Render
	AppController
	Db gorm.DB
}

func (c *PostsController) Show(rw http.ResponseWriter, r *http.Request) error {
	post := models.Post{}

	slug := mux.Vars(r)["slug"]
	c.Db.Where("slug = ?", slug).First(&post)

	c.HTML(rw, http.StatusOK, "posts/show", &models.PostView{Post: post})

	return nil
}
