package models

import (
	"gorm.io/gorm"
	"time"
)

type OrderMeal struct {
	Count         uint `json:"count" validate:"required"`
	CreateOrderId uint `json:"order_id"`
	MealId        uint `json:"meal_id"`
}

type CompleteOrderMeal struct {
	Count   uint `json:"count" validate:"required"`
	OrderId uint `json:"order_id"`
	Meal
}

type CreateOrder struct {
	Id           uint           `gorm:"primaryKey;autoIncrement:true" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	RestaurantId uint           `json:"restaurant_id" validate:"required"`
	UserId       uint           `json:"user_id" validate:"required"`
	Meals        []OrderMeal    `json:"meals" validate:"required"`
	Status       int            `json:"status" gorm:"default:0"`
}

type CompleteOrder struct {
	Id           uint                `gorm:"primaryKey;autoIncrement:true" json:"id"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
	DeletedAt    gorm.DeletedAt      `gorm:"index" json:"deleted_at"`
	RestaurantId uint                `json:"restaurant_id" validate:"required"`
	Restaurant   Restaurant          `json:"restaurant"`
	UserId       uint                `json:"user_id" validate:"required"`
	Meals        []CompleteOrderMeal `json:"meals" validate:"required"`
	Status       int                 `json:"status" gorm:"default:0"`
}

func NewCompleteOrder(order CreateOrder, meals []Meal, restaurant Restaurant) CompleteOrder {
	var completeMeals []CompleteOrderMeal
	for i, el := range order.Meals {
		if i < len(meals) {
			completeMeals = append(completeMeals, CompleteOrderMeal{
				Count:   el.Count,
				OrderId: el.CreateOrderId,
				Meal:    meals[i],
			})
		}
	}
	return CompleteOrder{
		Id:           order.Id,
		CreatedAt:    order.CreatedAt,
		UpdatedAt:    order.UpdatedAt,
		DeletedAt:    order.DeletedAt,
		RestaurantId: order.RestaurantId,
		UserId:       order.UserId,
		Meals:        completeMeals,
		Status:       order.Status,
		Restaurant:   restaurant,
	}
}
