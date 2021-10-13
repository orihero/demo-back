package models

import (
	"gorm.io/gorm"
	"time"
)

const (
	STATUS_PLACED = iota
	STATUS_CANCELED
	STATUS_PROCESSING
	STATUS_IN_ROUTE
	STATUS_DELIVERED
	STATUS_RECIEVED
)

type Restaurant struct {
	Id           uint           `gorm:"primaryKey;autoIncrement:true" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	UserId       uint           `json:"user_id" gorm:"not null" validate:"required"`
	Name         string         `json:"name" gorm:"not null" validate:"required"`
	Description  string         `json:"description" gorm:"not null" validate:"required"`
	PhotoUrl     string         `json:"photo_url"`
	BlockedUsers []User         `json:"blocked_users" gorm:"many2many:blocked_users"`
}

type Meal struct {
	Id           uint           `gorm:"primaryKey;autoIncrement:true" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	RestaurantId uint           `json:"restaurant_id" gorm:"not null" validate:"required"`
	Name         string         `json:"name" gorm:"not null" validate:"required"`
	Price        float64        `json:"price" gorm:"not null" validate:"required"`
	Description  string         `json:"description" gorm:"not null" validate:"required"`
	PhotoUrl     string         `json:"photo_url"`
}
