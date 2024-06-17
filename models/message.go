package models

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	Text        string `json:"text" gorm:"text;default:null"`
	SenderID    uint   `json:"-" gorm:"not null"`
	RecipientID uint   `json:"-" gorm:"not null"`
	Sender      User   `json:"sender" gorm:"foreignKey:SenderID"`
	Recipient   User   `json:"recipient" gorm:"foreignKey:RecipientID"`
}
