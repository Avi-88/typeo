package models

import "gorm.io/gorm"


type User struct {
	gorm.Model
	Username string `json:"username" gorm:"unique"`
	Email string `json:"email" gorm:"unique"`
	Password string `json:"password"`
}

type NewUser struct {
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
}