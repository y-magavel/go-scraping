package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func parseList(resp *http.Response) ([]Item, error) {
	body := resp.Body
	requestURL := *resp.Request.URL

	var items []Item

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, fmt.Errorf("get document error: %w", err)
	}

	tr := doc.Find("table tr")
	notFoundMessage := "ページが存在しません"
	if strings.Contains(doc.Text(), notFoundMessage) || tr.Size() == 0 {
		return nil, nil
	}

	tr.Each(func(_ int, s *goquery.Selection) {
		item := Item{}

		item.Name = s.Find("td:nth-of-type(2) a").Text()
		item.Price, _ = strconv.Atoi(strings.ReplaceAll(strings.ReplaceAll(s.Find("td:nth-of-type(3)").Text(), ",", ""), "円", ""))
		itemURL, exists := s.Find("td:nth-of-type(2) a").Attr("href")
		refURL, parseError := url.Parse(itemURL)

		if exists && parseError == nil {
			item.URL = (*requestURL.ResolveReference(refURL)).String()
		}

		if item.Name != "" {
			items = append(items, item)
		}
	})

	return items, nil
}
