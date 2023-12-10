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

func updateItemMaster(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var newItems []LatestItem
		err := tx.Unscoped().Joins("LEFT JOIN item_master ON latest_items.url = item_master.url").
			Where("item_master.url IS NULL").Find(&newItems).Error
		if err != nil {
			return fmt.Errorf("extract for bulk insert to item_master error: %w", err)
		}

		var insertRecords []ItemMaster
		for _, newItem := range newItems {
			insertRecords = append(insertRecords, ItemMaster{Item: newItem.Item})
			fmt.Printf("Index item is created: %s\n", newItem.URL)
		}
		if err := tx.CreateInBatches(insertRecords, 100).Error; err != nil {
			return fmt.Errorf("bulk insert to item_master error: %w", err)
		}

		var updatedItems []LatestItem
		if err := tx.Unscoped().Joins("INNER JOIN item_master ON latest_items.url = item_master.url").
			Where("latest_items.name <> item_master.name OR latest_items.price <> item_master.price OR item_master.deleted_at IS NOT NULL").
			Find(&updatedItems).Error; err != nil {
			return fmt.Errorf("update error: %w", err)
		}

		for _, updatedItem := range updatedItems {
			err := tx.Unscoped().Model(ItemMaster{}).Where("url = ?", updatedItem.URL).
				Updates(map[string]interface{}{"name": updatedItem.Name, "price": updatedItem.Price, "deleted_at": nil}).Error
			if err != nil {
				return fmt.Errorf("update error: %w", err)
			}
			fmt.Printf("Index item is updated: %s\n", updatedItem.URL)
		}

		var deletedItems []ItemMaster
		if err := tx.Where("NOT EXISTS(SELECT * FROM latest_items li WHERE li.url = item_master.url)").Find(&deletedItems).Error; err != nil {
			return fmt.Errorf("delete error: %w", err)
		}

		var ids []uint
		for _, deletedItem := range deletedItems {
			ids = append(ids, deletedItem.ID)
			fmt.Printf("Index item is deleted: %s\n", deletedItem.URL)
		}
		if len(ids) > 0 {
			if err := tx.Delete(&deletedItems).Error; err != nil {
				return fmt.Errorf("delete error: %w", err)
			}
		}

		return nil
	})
}
