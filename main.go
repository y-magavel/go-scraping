package main

import "fmt"

type Item struct {
	Name  string
	Price int
	URL   string
}

func main() {
	_, err := connectDB()

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
