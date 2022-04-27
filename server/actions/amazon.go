package actions

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
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

	respDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("RESPONSE:\n%s", string(respDump))

	json.NewDecoder(resp.Body).Decode(&results)

	return results
}
