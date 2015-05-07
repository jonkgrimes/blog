package controllers

import (
	"github.com/jinzhu/gorm"
	"gopkg.in/unrolled/render.v1"
)

type SessionsController struct {
	*render.Render
	AppController
	Db gorm.DB
}
