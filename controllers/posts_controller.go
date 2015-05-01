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

func (c *PostsController) Index(rw http.ResponseWriter, r *http.Request) error {
	var posts []models.Post

	c.Db.Find(&posts)

	c.JSON(rw, http.StatusOK, map[string]interface{}{"posts": posts})
	return nil
}

func (c *PostsController) Show(rw http.ResponseWriter, r *http.Request) error {
	post := models.Post{}

	id := mux.Vars(r)["id"]
	c.Db.First(&post, id)

	if post.Id != 0 {
		c.JSON(rw, http.StatusOK, map[string]interface{}{"post": post})
	} else {
		c.JSON(rw, http.StatusNotFound, nil)
	}

	return nil
}
