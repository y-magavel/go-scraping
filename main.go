package main

import (
	"os"
	"path/filepath"
	"reflect"
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
	Description         string
	ImageURL            string
	ImageLastModifiedAt time.Time
	ImageDownloadPath   string
	PDFURL              string
	PDFLastModifiedAt   time.Time
	PDFDownloadPath     string
}

func (i ItemMaster) TableName() string {
	return "item_master"
}

func (i ItemMaster) equals(target ItemMaster) bool {
	return reflect.DeepEqual(i, target)
}

func (i ItemMaster) ImageFileName() string {
	return filepath.Base(i.ImageURL)
}

func (i ItemMaster) PDFFileName() string {
	return filepath.Base(i.PDFURL)
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
	currentDir, _ := os.Getwd()
	downloadBasePath := filepath.Join(currentDir, "work", "downloadFiles")

	updatedItems, err = fetchDetails(updateChkItems, downloadBasePath)
	if err != nil {
		panic(err)
	}

	if err = createDetails(updatedItems, db); err != nil {
		panic(err)
	}
}
