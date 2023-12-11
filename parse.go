package main

import (
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
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

func parseDetail(response *http.Response, item ItemMaster, downloadBasePath string) (ItemMaster, error) {
	body := response.Body
	requestURL := *response.Request.URL
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return ItemMaster{}, fmt.Errorf("get detail page document body error %w", err)
	}

	item.Description = doc.Find("table tr:nth-of-type(2) td:nth-of-type(2)").Text()

	// Image
	href, exists := doc.Find("table tr:nth-of-type(1) td:nth-of-type(1) img").Attr("src")
	refURL, parseErr := url.Parse(href)
	if exists && parseErr == nil {
		imageURL := (*requestURL.ResolveReference(refURL)).String()
		isUpdated, currentLastModified := checkFileUpdated(imageURL, item.ImageLastModifiedAt)
		if isUpdated {
			item.ImageURL = imageURL
			item.ImageLastModifiedAt = currentLastModified

			imageDownloadPath := filepath.Join(downloadBasePath, "img", strconv.Itoa(int(item.ID)), item.ImageFileName())
			err := downloadFile(imageURL, imageDownloadPath)
			if err != nil {
				return ItemMaster{}, fmt.Errorf("download image error: %w", err)
			}
			item.ImageDownloadPath = imageDownloadPath
		}
	}

	// PDF
	href, exists = doc.Find("table tr:nth-of-type(3) td:nth-of-type(2) a").Attr("href")
	refURL, parseErr = url.Parse(href)
	if exists && parseErr == nil {
		pdfURL := (*requestURL.ResolveReference(refURL)).String()
		isUpdated, currentLastModified := checkFileUpdated(pdfURL, item.PDFLastModifiedAt)
		if isUpdated {
			item.PDFURL = pdfURL
			item.PDFLastModifiedAt = currentLastModified

			pdfDownloadPath := filepath.Join(downloadBasePath, "pdf", strconv.Itoa(int(item.ID)), item.PDFFileName())
			err := downloadFile(pdfURL, pdfDownloadPath)
			if err != nil {
				return ItemMaster{}, fmt.Errorf("download pdf error: %w", err)
			}
			item.PDFDownloadPath = pdfDownloadPath
		}
	}

	return item, nil
}
