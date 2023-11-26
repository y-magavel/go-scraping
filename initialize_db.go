package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func connectDB() (*gorm.DB, error) {
	dbHost := "localhost"
	dbPort := "3306"
	dbName := "go_scraping_dev"
	dbUser := "go-scraping-user"
	dbPassword := "password"

	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Errorf("db connection error: %w", err)
	}

	return db, nil
}
