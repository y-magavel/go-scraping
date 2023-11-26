package main

import (
	"fmt"

	"gorm.io/gorm"
)

func createLatestItems(items []Item, db *gorm.DB) error {
	stmt := &gorm.Statement{DB: db}
	err := stmt.Parse(&LatestItem{})
	if err != nil {
		fmt.Errorf("get latest_items table name error: %w", err)
	}

	if err := db.Exec("TRUNCATE " + stmt.Schema.Table).Error; err != nil {
		fmt.Errorf("truncate latest_items error: %w", err)
	}

	var insertRecords []LatestItem
	for _, item := range items {
		insertRecords = append(insertRecords, LatestItem{Item: item})
	}

	if err := db.CreateInBatches(insertRecords, 100).Error; err != nil {
		return fmt.Errorf("bulk insert to latest_items error: %w", err)
	}

	return nil
}
