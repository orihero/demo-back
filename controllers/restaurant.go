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

func CreateRestaurant(w http.ResponseWriter, r *http.Request) {
	//Only admins are allowed to use this method
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
	//Verify the owner of a restaurant
	email := r.Header.Get("Email")
	user := models.User{Email: email}
	connection.Where(&user).First(&user)

	var restaurant models.Restaurant
	err := json.NewDecoder(r.Body).Decode(&restaurant)
	if err != nil {
		var err models.Error
		err = utils.SetError(err, "Unprocessable entity.")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}
	//Bind restaurant to the current user
	restaurant.UserId = user.Id
	//Validation
	v := validator.New()
	if errs := v.Struct(restaurant); errs != nil {
		var err models.Error
		err = utils.SetError(err, "Validation error:\n"+fmt.Sprint(errs))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}
	//Success
	connection.Create(&restaurant)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(restaurant)
}

func GetUserRestaurants(w http.ResponseWriter, r *http.Request) {
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
	var restaurants []models.Restaurant
	fmt.Println(user.Id, user.Email, email)
	connection.Where(&models.Restaurant{UserId: user.Id}).Find(&restaurants)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(restaurants)
}

func GetRestaurantMeals(w http.ResponseWriter, r *http.Request) {
	connection := database.GetDatabase()
	defer database.CloseDatabase(connection)
	restaurantId, err := strconv.ParseUint(r.URL.Query().Get("restaurant_id"), 10, 32)
	if err != nil {
		var err models.Error
		err = utils.SetError(err, "Invalid restaurant ID")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}
	meal := models.Meal{RestaurantId: uint(restaurantId)}
	var meals []models.Meal
	connection.Where(&meal).Find(&meals)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(meals)
}

func GetAllMeals(w http.ResponseWriter, r *http.Request) {
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
	var restaurants []models.Restaurant
	fmt.Println(user.Id, user.Email, email)
	connection.Where(&models.Restaurant{UserId: user.Id}).Find(&restaurants)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(restaurants)
}

func GetAllRestaurants(w http.ResponseWriter, r *http.Request) {
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
	var restaurants []models.Restaurant
	connection.Preload("BlockedUsers").Find(&restaurants)
	w.Header().Set("Content-Type", "text/html; charset=ascii")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers","Content-Type,access-control-allow-origin, access-control-allow-headers")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(restaurants)
}

func CreateMeal(w http.ResponseWriter, r *http.Request) {
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
	var meal models.Meal
	err := json.NewDecoder(r.Body).Decode(&meal)
	v := validator.New()
	if errs := v.Struct(meal); errs != nil {
		var err models.Error
		err = utils.SetError(err, "Validation error:\n"+fmt.Sprint(errs))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}
	fmt.Println(meal)
	if err != nil {
		var err models.Error
		err = utils.SetError(err, "Unprocessable entity.")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}
	connection.Create(&meal)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(meal)
}

func _isAdmin(r *http.Request) bool {
	role := r.Header.Get("Role")
	return role == "admin"
}
