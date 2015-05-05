package models

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id                int64  `json:"id"`
	Username          string `sql:"size:127;unique_index" json:"username"`
	Email             string `sql:"unique_index" json:"email"`
	Password          string `sql:"-" validate:"required,max=60"`
	EncryptedPassword string `json:"-"`
}

func (u *User) BeforeSave() error {
	u.encryptPassword()

	return nil
}

func (u *User) encryptPassword() error {
	password := []byte(u.Password)
	encryptedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	u.EncryptedPassword = string(encryptedPassword)

	return err
}
