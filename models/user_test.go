package models

import (
	"golang.org/x/crypto/bcrypt"
	"testing"
)

var encryptPasswordCases = []struct {
	user     User
	password []byte
}{
	{
		User{Password: "password"},
		[]byte("password"),
	},
	{
		User{Password: "Haxx0r$"},
		[]byte("Haxx0r$"),
	},
}

func TestGenerateEncryptPassword(t *testing.T) {
	t.Log("TestGenerateEncryptedPassword")

	for _, testCase := range encryptPasswordCases {
		user := testCase.user
		user.encryptPassword()

		matched := bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), testCase.password) == nil

		if !matched {
			t.Logf("The hashed password did not contain the password")
			t.Fail()
		}
	}
}
