package models

import (
	"eda/utils/token"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/html"
	"gorm.io/gorm"
	"strings"
)

const RoleAdmin = "admin"
const RoleModer = "moder"
const RoleGuest = "guest"

type User struct {
	gorm.Model
	Phone    string `json:"phone" gorm:"text;unique"`
	Email    string `json:"email" gorm:"text;unique"`
	Password string `json:"password" gorm:"text;size:255"`
	Role     string `json:"role" gorm:"text;not null;default:'guest'"`
	IsActive bool   `json:"is_active" gorm:"bool;not null;default:false"`
}

func GetUserByID(uid uint) (User, error) {
	var u User

	if err := DB.First(&u, uid).Error; err != nil {
		return u, errors.New("user not found")
	}

	u.PrepareGive()

	return u, nil
}

func GetUserByEmail(email string) (User, error) {
	var u User
	err := DB.Model(User{}).Where("email = ?", email).Find(&u).Error
	if err != nil {
		return u, errors.New("user not found")
	}

	return u, nil
}

func GetUserByPhone(phone string) (User, error) {
	var u User
	err := DB.Model(User{}).Where("phone = ?", phone).Find(&u).Error
	if err != nil {
		return u, errors.New("user not found")
	}

	return u, nil
}

func GetUserRoleById(uid uint) (string, error) {
	var u User

	if err := DB.First(&u, uid).Error; err != nil {
		return u.Role, errors.New("user not found")
	}

	return u.Role, nil
}

func (u *User) PrepareGive() {
	u.Password = ""
}

func (u *User) SaveUser() (User, error) {
	var err error
	err = DB.Create(&u).Error
	if err != nil {
		return User{}, err
	}
	return *u, nil
}

func (u *User) SetActive() (User, error) {
	var err error
	err = DB.Model(User{}).Update("is_active", true).Error
	if err != nil {
		return User{}, err
	}
	return *u, nil
}

func (u *User) BeforeSave(_ *gorm.DB) error {
	//turn password into hash
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)

	//remove spaces in email
	u.Phone = html.EscapeString(strings.TrimSpace(u.Phone))

	return nil
}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func LoginCheck(Phone string, password string) (string, error) {
	var err error

	u := User{}

	err = DB.Model(User{}).Where("phone = ?", Phone).Take(&u).Error

	if err != nil {
		return "", err
	}

	err = VerifyPassword(password, u.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}

	generatedToken, err := token.GenerateToken(u.ID)

	if err != nil {
		return "", err
	}

	return generatedToken, nil
}
