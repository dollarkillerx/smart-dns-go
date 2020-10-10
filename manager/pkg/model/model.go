package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `json:"username"`
}

type Domain struct {
	
}

type Route struct {
	
}