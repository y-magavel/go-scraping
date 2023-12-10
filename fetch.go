package main

import (
	"fmt"
	"net/http"
)

func fetch(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http get request error: %w", err)
	}

	return resp, nil
}

func fetchDetails(items []ItemMaster) ([]ItemMaster, error) {
	var updatedItems []ItemMaster

	for _, item := range items {
		response, err := fetch(item.URL)
		if err != nil {
			return nil, fmt.Errorf("fetch detail page body error: %w", err)
		}

		currentItem, err := parseDetail(response, item)
		if err != nil {
			return nil, fmt.Errorf("fetch detail page content error: %w", err)
		}

		if !item.equals(currentItem) {
			updatedItems = append(updatedItems, currentItem)
		}
	}

	return updatedItems, nil
}
