package actions

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"golang.org/x/net/html"
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
	var results []AmazonSearchResultsPage

	str := strings.Join(strings.Split(keyword, " "), "+")

	serp := fmt.Sprintf("https://www.amazon.com/s?k=%s", str)

	host := os.Getenv("P_HOST")
	username := os.Getenv("P_USERNAME")
	sessionId := fmt.Sprint(rand.Intn(1000000))
	path := username + sessionId + ":" + host

	u, err := url.Parse(path)
	if err != nil {
		log.Fatal(err)
		return results
	}

	tr := &http.Transport{
		Proxy: http.ProxyURL(u),
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client := &http.Client{
		Transport: tr,
	}

	req, err := http.NewRequest("GET", serp, nil)

	if err != nil {
		fmt.Println("Request failed: ", err)
		return results
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error while fetching Amazon SERP", err)
		return results
	}
	defer resp.Body.Close()

	ParseHtml(resp.Body)

	if err != nil {
		fmt.Println("Error while writing response body from Amazon crawl.")
		return results
	}

	json.NewDecoder(resp.Body).Decode(&results)

	return results
}

func ParseHtml(r io.Reader) (string, error) {

	reviewsRegex := regexp.MustCompile("[0-9,]+")
	moneyRegex := regexp.MustCompile(`^\/[A-Z]([A-Za-z0-9]*[-%\/])+(\/dp\/[A-Z]+)*`)
	amazonASIN := regexp.MustCompile(`(\/[A-Z0-9]{10,}\/)`)

	doc := html.NewTokenizer(r)

	for {
		el := doc.Next()
		var product AmazonSearchResultsPage

		switch el {
		case html.ErrorToken:
			return "", doc.Err()
		case html.StartTagToken, html.EndTagToken:
			e := doc.Token()
			if strings.Contains(e.String(), "a href=") {
				el = doc.Next()
				if el == html.TextToken {
					if amazonASIN.MatchString(doc.Token().Data) {
						product.Name = doc.Token().Data
						fmt.Println(doc.Token().Data)
					}
				}
			}
			if strings.Contains(e.String(), "a-size-base") {
				el = doc.Next()
				if el == html.TextToken {
					if reviewsRegex.MatchString(doc.Token().Data) {
						product.Reviews = doc.Token().Data
						fmt.Println(doc.Token().Data)
					}
					if moneyRegex.MatchString(doc.Token().Data) {
						product.Price = doc.Token().Data
					}
				}
			}
		}
	}
}
