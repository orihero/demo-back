package controllers

import (
	"../database"
	"../models"
	"../utils"
	"encoding/json"
	"fmt"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"strconv"
)

func RestaurantOrders(w http.ResponseWriter, r *http.Request) {
	if !_isAdmin(r) {
		var err models.Error
		err = utils.SetError(err, "Method not  allowed.")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(err)
		return
	}
	connection := database.GetDatabase()
	defer database.CloseDatabase(connection)
	email := r.Header.Get("Email")
	user := models.User{Email: email}
	connection.Where(&user).First(&user)
	query := models.Restaurant{UserId: user.Id}
	var userRestaurants []models.Restaurant
	var restaurantsOrders []models.CompleteOrder
	connection.Where(&query).Find(&userRestaurants)
	for _,el:=range userRestaurants{
		var orders []models.CreateOrder
		connection.Where(&models.CreateOrder{RestaurantId: el.Id}).Order("created_at desc").Preload("Meals").Find(&orders)
		for _, order := range orders {
			var ids []uint
			for _, meal := range order.Meals {
				ids = append(ids, meal.MealId)
			}
			var meals []models.Meal
			connection.Find(&meals, ids)
			restaurant := models.Restaurant{Id: order.RestaurantId}
			connection.Where(&restaurant).Find(&restaurant)
			restaurantsOrders = append(restaurantsOrders, models.NewCompleteOrder(order, meals, restaurant))
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(restaurantsOrders)
}

func UserOrders(w http.ResponseWriter, r *http.Request) {
	if _isAdmin(r) {
		var err models.Error
		err = utils.SetError(err, "Method not  allowed.")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(err)
		return
	}
	connection := database.GetDatabase()
	defer database.CloseDatabase(connection)
	email := r.Header.Get("Email")
	user := models.User{Email: email}
	connection.Where(&user).First(&user)
	order := models.CreateOrder{UserId: user.Id}
	var orders []models.CreateOrder
	connection.Where(&order).Order("created_at desc").Preload("Meals").Find(&orders)
	var completeOrders []models.CompleteOrder
	for _, el := range orders {
		var ids []uint
		for _, meal := range el.Meals {
			ids = append(ids, meal.MealId)
		}
		var meals []models.Meal
		connection.Find(&meals, ids)
		restaurant := models.Restaurant{Id: el.RestaurantId}
		connection.Where(&restaurant).Find(&restaurant)
		completeOrders = append(completeOrders, models.NewCompleteOrder(el, meals, restaurant))
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(completeOrders)
}

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	if _isAdmin(r) {
		var err models.Error
		err = utils.SetError(err, "Method not  allowed.")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(err)
		return
	}
	connection := database.GetDatabase()
	defer database.CloseDatabase(connection)
	var order models.CreateOrder
	email := r.Header.Get("Email")
	user := models.User{Email: email}
	connection.Where(&user).First(&user)
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		var err models.Error
		err = utils.SetError(err, "Unprocessable entity.")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err)
		return
	}
	//TODO check if meals belong to the current restaurant
	order.UserId = user.Id
	order.Status = models.STATUS_PLACED
	v := validator.New()
	if errs := v.Struct(order); errs != nil {
		var err models.Error
		err = utils.SetError(err, "Validation error:\n"+fmt.Sprint(errs))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}
	connection.Create(&order)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func ChangeOrderStatus(w http.ResponseWriter, r *http.Request) {

	connection := database.GetDatabase()
	defer database.CloseDatabase(connection)
	orderId, err := strconv.ParseUint(r.URL.Query().Get("order_id"), 10, 32)
	//CreateOrder id invalid
	if err != nil {
		var err models.Error
		err = utils.SetError(err, "Invalid order ID")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}
	status, err := strconv.Atoi(r.URL.Query().Get("status"))
	//Invalid status
	if err != nil {
		var err models.Error
		err = utils.SetError(err, "Invalid status")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}
	//Get order
	order := models.CreateOrder{Id: uint(orderId)}
	connection.Where(&order).Preload("Meals").First(&order)
	//Vars for user roles
	role := r.Header.Get("Role")
	//Invalid order status message generator
	invalidOrderStatus := func() {
		var err models.Error
		err = utils.SetError(err, "Invalid order status.")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
	}
	notAllowedStatusChangeAttempt := func() {
		var err models.Error
		err = utils.SetError(err, "Not allowed status change attempt")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(err)
		return
	}
	//Next status
	switch status {
	case models.STATUS_CANCELED, models.STATUS_RECIEVED:
		if role != models.ROLE_USER {
			notAllowedStatusChangeAttempt()
		}
		order.Status = status
		connection.Save(&order)
		break
	case models.STATUS_PROCESSING, models.STATUS_IN_ROUTE, models.STATUS_DELIVERED:
		if role != models.ROLE_ADMIN {
			notAllowedStatusChangeAttempt()
			return
		}
		order.Status = status
		connection.Save(&order)
	default:
		invalidOrderStatus()
		return
	}
}
