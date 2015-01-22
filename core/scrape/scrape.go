package scrape

import (
	"github.com/PuerkitoBio/goquery"
)

type ScrapeRequest interface {
	Scrape(doc *goquery.Document) map[string]interface{}
}
