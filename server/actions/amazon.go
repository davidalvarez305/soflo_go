package actions

import (
	"crypto/tls"
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

	"github.com/PuerkitoBio/goquery"
)

type AmazonSearchResultsPage struct {
	Image    string `json:"image"`
	Name     string `json:"name"`
	Link     string `json:"link"`
	Reviews  string `json:"reviews"`
	Price    string `json:"price"`
	Rating   string `json:"rating"`
	Category string `json:"category"`
}

func FetchCrawler(keyword string, category string) []AmazonSearchResultsPage {
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

	products, err := ParseHtml(resp.Body, category)
	results = products

	if err != nil {
		fmt.Println("Error while parsing HTML.")
	}

	return results
}

func ParseHtml(r io.Reader, category string) ([]AmazonSearchResultsPage, error) {
	var products []AmazonSearchResultsPage

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		fmt.Println("Error trying to parse document.")
		return products, err
	}

	doc.Find(".sg-col-inner").Each(func(i int, s *goquery.Selection) {
		var product AmazonSearchResultsPage

		reviewsRegex := regexp.MustCompile("[0-9,]+")
		moneyRegex := regexp.MustCompile(`[\$]+?(\d+([,\.\d]+)?)`)
		amazonASIN := regexp.MustCompile(`(\/[A-Z0-9]{10,}\/)`)

		el, _ := s.Find("a").Attr("href")
		cond := amazonASIN.MatchString(el)

		if cond {
			name := strings.Join(strings.Split(strings.Split(el, "/")[1], "-"), " ")
			product.Name = name

			rating := strings.Split(s.Find(".a-icon-alt").Text(), " ")[0]
			product.Rating = rating

			link := strings.Split(el, "/")[3]
			product.Link = "https://amazon.com/dp/" + link + "?tag=sfac09-20&linkCode=ogi&th=1&psc=1"

			image, _ := s.Find("img").Attr("src")
			product.Image = image

			product.Category = category

			if len(moneyRegex.FindAllString(s.Find(".a-size-base").Text(), 3)) > 0 {
				price := moneyRegex.FindAllString(s.Find(".a-size-base").Text(), 3)[0]
				product.Price = price
			}
			if len(reviewsRegex.FindAllString(s.Find(".a-size-base").Text(), 2)) > 0 {
				reviews := reviewsRegex.FindAllString(s.Find(".a-size-base").Text(), 3)[0]
				product.Reviews = reviews

			}
			products = append(products, product)
		}
	})
	return products, nil
}
