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

	id := mux.Vars(r)["id"]
	c.Db.First(&post, id)

	c.HTML(rw, http.StatusOK, "posts/show", &models.PostView{Post: post})

	return nil
}
