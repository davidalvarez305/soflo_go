package actions

import (
	"fmt"
	"strings"
	"time"
)

type AmazonSearchResultsPage struct {
	Image    string `json:"image"`
	Name     string `json:"name"`
	Link     string `json:"link"`
	Reviews  string `json:"reviews"`
	Price    string `json:"price"`
	Rating   string `json:"rating"`
	ID       string `json:"id"`
	Category string `json:"category"`
}

func FetchCrawler(keyword string) []AmazonSearchResultsPage {
	time.Sleep(1 * time.Second)
	var data []AmazonSearchResultsPage

	str := strings.Join(strings.Split(keyword, " "), "+")

	serp := fmt.Sprintf("https://www.amazon.com/s?k=%s", str)

	fmt.Println(serp)

	return data
}
