package main

import (
	"fmt"
	"net/http"
	"time"
)

func fetch(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http get request error: %w", err)
	}

	return resp, nil
}

func fetchDetails(items []ItemMaster, downloadBasePath string) ([]ItemMaster, error) {
	var updatedItems []ItemMaster

	for _, item := range items {
		response, err := fetch(item.URL)
		if err != nil {
			return nil, fmt.Errorf("fetch detail page body error: %w", err)
		}

		currentItem, err := parseDetail(response, item, downloadBasePath)
		if err != nil {
			return nil, fmt.Errorf("fetch detail page content error: %w", err)
		}

		if !item.equals(currentItem) {
			updatedItems = append(updatedItems, currentItem)
		}
	}

	return updatedItems, nil
}

func checkFileUpdated(fileURL string, lastModified time.Time) (isUpdated bool, currentLastModified time.Time) {
	getLastModified := func(fileURL string) (time.Time, error) {
		res, err := http.Head(fileURL)
		if err != nil {
			return time.Time{}, fmt.Errorf("http head request error: %w", err)
		}
		lastModified, err := time.Parse("Mon, 02 Jan 2006 15:04:05 MST", res.Header.Get("Last-Modified"))
		if err != nil {
			return time.Time{}, fmt.Errorf("get last-modified attribute error: %w", err)
		}
		return lastModified, nil
	}

	currentLastModified, err := getLastModified(fileURL)
	if err != nil {
		return false, currentLastModified
	}

	if currentLastModified.After(lastModified) {
		return true, currentLastModified
	} else {
		return false, lastModified
	}
}
