package model

import (
	"gorm.io/gorm"
	"time"
)

type Model struct {
	ID        uint           `gorm:"primarykey" json:"-"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type User struct {
	Model
	Email    string `gorm:"column:email;index" json:"email"`
	Username string `gorm:"column:username;index" json:"username"`
	Password string `gorm:"column:username;index" json:"password"`
	Salt     string `gorm:"column:salt" json:"-"`
}

type Domain struct {
}

type Route struct {
}
