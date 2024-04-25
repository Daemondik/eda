package models

import (
	"eda/database"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/html"
	"gorm.io/gorm"
	"strings"
)

type Fact struct {
	gorm.Model
	Question string `json:"question" gorm:"text;not null;default:null"`
	Answer   string `json:"answer" gorm:"text;not null;default:null"`
}

type User struct {
	gorm.Model
	Email    string `json:"email" gorm:"text;not null;unique"`
	Password string `json:"password" gorm:"text;size:255;not null;"`
}

func (u *User) SaveUser() (*User, error) {
	var err error
	err = database.DB.Db.Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) BeforeSave() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))

	return nil
}
