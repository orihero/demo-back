package database

import (
	"../models"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

//-------------DATABASE FUNCTIONS---------------------

//returns database connection
func GetDatabase() *gorm.DB {
	connection, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		log.Fatalln("Invalid database url")
	}
	sql, err := connection.DB()

	err = sql.Ping()
	if err != nil {
		log.Fatal("Database connected")
	}
	fmt.Println("Database connection successful.")
	return connection
}

//create user table in userdb
func InitialMigration() {
	connection := GetDatabase()
	defer CloseDatabase(connection)
	_ = connection.AutoMigrate(models.User{})
	_ = connection.AutoMigrate(models.Restaurant{})
	_ = connection.AutoMigrate(models.Meal{})
	_ = connection.AutoMigrate(models.CreateOrder{})
	_ = connection.AutoMigrate(models.OrderMeal{})
}

//closes database connection
func CloseDatabase(connection *gorm.DB) {
	sqldb, _ := connection.DB()
	sqldb.Close()
}
