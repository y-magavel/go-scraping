package main

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Item struct {
	Name  string `gorm:"type:varchar(100);not null;"`
	Price int
	URL   string `gorm:"type:varchar(100);uniqueIndex;"`
}

type LatestItem struct {
	Item
	CreatedAt time.Time
}

type ItemMaster struct {
	gorm.Model
	Item
}

func (ItemMaster) TableName() string {
	return "item_master"
}

func main() {
	db, err := connectDB()
	if err != nil {
		panic(err)
	}

	err = migrateDB(db)
	if err != nil {
		panic(err)
	}

	baseURL := "http://localhost:5001"
	resp, err := fetch(baseURL)
	if err != nil {
		panic(err)
	}

	indexItems, err := parseList(resp)
	if err != nil {
		panic(err)
	}

	for _, item := range indexItems {
		fmt.Println(item)
	}
}
