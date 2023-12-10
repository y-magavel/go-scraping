package main

import (
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
	Description string
}

func (i ItemMaster) TableName() string {
	return "item_master"
}

func (i ItemMaster) equals(target ItemMaster) bool {
	return i.Description == target.Description
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

	items, err := parseList(resp)
	if err != nil {
		panic(err)
	}

	if err := createLatestItems(items, db); err != nil {
		panic(err)
	}

	if err := updateItemMaster(db); err != nil {
		panic(err)
	}

	var updateChkItems []ItemMaster
	updateChkItems, err = findItemMaster(db)
	if err != nil {
		panic(err)
	}

	var updatedItems []ItemMaster
	updatedItems, err = fetchDetails(updateChkItems)
	if err != nil {
		panic(err)
	}

	if err = createDetails(updatedItems, db); err != nil {
		panic(err)
	}
}
