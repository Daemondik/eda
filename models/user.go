package models

import (
	"eda/utils/token"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"strings"
)

const RoleAdmin = "admin"
const RoleModer = "moder"
const RoleGuest = "guest"

type User struct {
	gorm.Model
	Phone    string `json:"phone" gorm:"text"`
	Email    string `json:"email" gorm:"text"`
	Password string `json:"password" gorm:"text;size:255"`
	Role     string `json:"role" gorm:"text;not null;default:'guest'"`
	IsActive bool   `json:"is_active" gorm:"bool;not null;default:false"`
}

func GetUserByID(uid uint) (User, error) {
	var u User
	if err := DB.Where("is_active = true").First(&u, uid).Error; err != nil {
		return User{}, errors.New("user not found")
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
	err := DB.Where("phone = ?", phone).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return u, errors.New("user not found")
		}
		return u, err
	}

	return u, nil
}

func GetUserRoleById(uid uint) (string, error) {
	var u User
	if err := DB.First(&u, uid).Error; err != nil {
		return "", errors.New("user not found")
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
	err := DB.Model(u).Where("id = ?", u.ID).Update("is_active", true).Error
	if err != nil {
		return User{}, err
	}
	return *u, nil
}

func (u *User) BeforeSave(_ *gorm.DB) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	u.Phone = strings.TrimSpace(u.Phone)
	return nil
}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func LoginCheck(phone string, password string) (string, error) {
	var err error
	u := User{}
	err = DB.Model(User{}).Where("phone = ?", phone).Take(&u).Error
	if err != nil {
		return "", err
	}
	err = VerifyPassword(password, u.Password)
	if err != nil && errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return "", err
	}
	generatedToken, err := token.GenerateToken(u.ID)
	if err != nil {
		return "", err
	}
	return generatedToken, nil
}
