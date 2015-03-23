package controllers

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"gopkg.in/unrolled/render.v1"

	"blog/models"
)

type HomeController struct {
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
