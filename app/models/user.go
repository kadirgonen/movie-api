package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Id        string `gorm:"unique"`
	FirstName string
	LastName  string
	Email     string
	Password  string
	IsAdmin   bool `gorm:"default:false"`
}
