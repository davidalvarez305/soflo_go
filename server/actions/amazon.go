package actions

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	strftime "github.com/itchyny/timefmt-go"

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

type AmazonPaapi5RequestBody struct {
	Marketplace string   `json:"Marketplace"`
	PartnerType string   `json:"PartnerType"`
	PartnerTag  string   `json:"PartnerTag"`
	Keywords    string   `json:"Keywords"`
	SearchIndex string   `json:"SearchIndex"`
	ItemCount   int      `json:"ItemCount"`
	Resources   []string `json:"Resources"`
}

func ScrapeSearchResultsPage(keyword string) []AmazonSearchResultsPage {
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

	products, err := parseHtml(resp.Body, keyword)
	results = products

	if err != nil {
		fmt.Println("Error while parsing HTML.")
	}

	return results
}

func parseHtml(r io.Reader, keyword string) ([]AmazonSearchResultsPage, error) {
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
			product.Link = "https://amazon.com/dp/" + link + os.Getenv("AMAZON_TAG")

			image, _ := s.Find("img").Attr("src")
			product.Image = image

			product.Category = keyword

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

func makeHash(hash hash.Hash, b []byte) []byte {
	hash.Reset()
	hash.Write(b)
	return hash.Sum(nil)
}

func buildCanonicalString(method, uri, query, signedHeaders, canonicalHeaders, payloadHash string) string {
	return strings.Join([]string{
		method,
		uri,
		query,
		canonicalHeaders,
		signedHeaders,
		payloadHash,
	}, "\n")
}

func buildStringToSign(date, credentialScope, canonicalRequest string) string {
	return strings.Join([]string{
		"AWS4-HMAC-SHA256",
		date,
		credentialScope,
		hex.EncodeToString(makeHash(sha256.New(), []byte(canonicalRequest))),
	}, "\n")
}

func buildCanonicalHeaders(contentEncoding, host, xAmazonDate, amazonTarget string) string {
	return strings.Join([]string{"content-encoding:" + contentEncoding, "host:" + host, "x-amz-date:" + xAmazonDate, "x-amz-target" + amazonTarget}, "\n")
}

func buildSignature(strToSign string, sig string) (string, error) {
	return hex.EncodeToString(HMACSHA256([]byte(sig), []byte(strToSign))), nil
}

func HMACSHA256(key []byte, data []byte) []byte {
	hash := hmac.New(sha256.New, key)
	hash.Write(data)
	return hash.Sum(nil)
}

func SearchPaapi5Items(keyword string) []AmazonSearchResultsPage {
	var products []AmazonSearchResultsPage

	resources := []string{
		"Images.Primary.Medium",
		"ItemInfo.Title",
		"Offers.Listings.Price",
		"ItemInfo.ByLineInfo",
		"ItemInfo.Features",
		"ItemInfo.ProductInfo"}

	d := AmazonPaapi5RequestBody{
		Marketplace: "www.amazon.com",
		PartnerType: "Associates",
		PartnerTag:  os.Getenv("AMAZON_PARTNER_TAG"),
		Keywords:    keyword,
		SearchIndex: "All",
		ItemCount:   10,
		Resources:   resources,
	}

	body, e := json.Marshal(d)

	if e != nil {
		return products
	}

	method := "POST"
	service := "ProductAdvertisingAPI"
	url := "https://webservices.amazon.com/paapi5/searchitems"
	host := "webservices.amazon.com"
	region := os.Getenv("AWS_REGION")
	contentType := "application/json; charset=UTF-8"
	amazonTarget := "com.amazon.paapi5.v1.ProductAdvertisingAPIv1.SearchItems"
	contentEncoding := "amz-1.0"
	t := time.Now()
	amazonDate := strftime.Format(t, "%Y%m%d")
	xAmazonDate := strftime.Format(t, "%Y%m%dT%H%M%SZ")
	canonicalUri := "/paapi5/searchitems"
	canonicalQuerystring := ""
	canonicalHeaders := buildCanonicalHeaders(contentEncoding, host, xAmazonDate, amazonTarget)
	credentialScope := amazonDate + "/" + region + "/" + service + "/" + "aws4_request"
	signedHeaders := "content-encoding;host;x-amz-date;x-amz-target"

	kSecret := os.Getenv("AWS_SECRET_ACCESS_KEY")
	kDate := hex.EncodeToString(HMACSHA256([]byte("AWS4"+kSecret), []byte(amazonDate)))
	kRegion := hex.EncodeToString(HMACSHA256([]byte(kDate), []byte(region)))
	kService := hex.EncodeToString(HMACSHA256([]byte(kRegion), []byte(service)))
	signingKey := hex.EncodeToString(HMACSHA256([]byte(kService), []byte("aws4_request")))

	canonicalRequest := buildCanonicalString(method, canonicalUri, canonicalQuerystring, signedHeaders, canonicalHeaders, hex.EncodeToString(makeHash(sha256.New(), body)))
	stringToSign := buildStringToSign(xAmazonDate, credentialScope, canonicalRequest)
	signature, err := buildSignature(stringToSign, signingKey)
	if err != nil {
		fmt.Println("Error while building signature.")
	}

	authorizationHeader := "AWS4-HMAC-SHA256" + " Credential=" + os.Getenv("AWS_ACCESS_KEY_ID") + "/" + credentialScope + ", SignedHeaders=" + signedHeaders + ", Signature=" + signature
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))

	if err != nil {
		fmt.Println("Request failed: ", err)
		return products
	}

	req.Header.Set("content-type", contentType)
	req.Header.Set("content-encoding", contentEncoding)
	req.Header.Set("host", host)
	req.Header.Set("x-amz-date", xAmazonDate)
	req.Header.Set("x-amz-target", amazonTarget)
	req.Header.Set("Authorization", authorizationHeader)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error while fetching Amazon SERP", err)
		return products
	}
	defer resp.Body.Close()

	respDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("RESPONSE:\n%s", string(respDump))

	return products
}
