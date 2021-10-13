package models

import (
	"gorm.io/gorm"
	"time"
)

const (
	ROLE_ADMIN = "admin"
	ROLE_USER  = "user"
)

//------------Auth models---------------------

type User struct {
	Id          uint           `gorm:"primaryKey;autoIncrement:true" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Name        string         `json:"name"`
	Email       string         `gorm:"unique" json:"email"`
	Password    string         `json:"password"`
	Role        string         `json:"role"`
	Restaurants []Restaurant   `json:"restaurants"`
	IsBlocked   bool           `json:"is_blocked"`
}

type Authentication struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
