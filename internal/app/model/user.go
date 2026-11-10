package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	Base
	Username    string `json:"username" gorm:"unique"`
	Password    string `json:"-"`
	Email       string `json:"email" gorm:"unique"`
	DisplayName string `json:"display_name"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return
}
