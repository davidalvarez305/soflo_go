package actions

import (
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/davidalvarez305/soflo_go/server/types"
	"github.com/davidalvarez305/soflo_go/server/utils"
)

func ScrapeSearchResultsPage(keyword string) []types.AmazonSearchResultsPage {
	var results []types.AmazonSearchResultsPage

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

	products, err := utils.ParseHtml(resp.Body, keyword)
	results = products

	if err != nil {
		fmt.Println("Error while parsing HTML.")
		return products
	}

	return results
}

func SearchPaapi5Items(keyword string) types.PAAAPI5Response {
	var products types.PAAAPI5Response

	resources := []string{
		"Images.Primary.Medium",
		"ItemInfo.Title",
		"Offers.Listings.Price",
		"ItemInfo.ByLineInfo",
		"ItemInfo.Features",
		"ItemInfo.ProductInfo"}

	d := types.AmazonPaapi5RequestBody{
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
	amazonDate := utils.FormatShortDate(t)
	xAmazonDate := utils.FormatDate(t)
	canonicalUri := "/paapi5/searchitems"
	canonicalQuerystring := ""
	canonicalHeaders := utils.BuildCanonicalHeaders(contentType, contentEncoding, host, xAmazonDate, amazonTarget)
	credentialScope := amazonDate + "/" + region + "/" + service + "/" + "aws4_request"
	signedHeaders := "content-encoding;content-type;host;x-amz-date;x-amz-target"

	kSecret := os.Getenv("AWS_SECRET_ACCESS_KEY")
	kDate := utils.HMACSHA256([]byte("AWS4"+kSecret), []byte(amazonDate))
	kRegion := utils.HMACSHA256(kDate, []byte(region))
	kService := utils.HMACSHA256(kRegion, []byte(service))
	signingKey := utils.HMACSHA256(kService, []byte("aws4_request"))

	canonicalRequest := utils.BuildCanonicalString(method, canonicalUri, canonicalQuerystring, signedHeaders, canonicalHeaders, hex.EncodeToString(utils.MakeHash(sha256.New(), body)))
	stringToSign := utils.BuildStringToSign(xAmazonDate, credentialScope, canonicalRequest)
	signature, err := utils.BuildSignature(stringToSign, signingKey)
	if err != nil {
		fmt.Println("Error while building signature.")
		return products
	}

	authorizationHeader := "AWS4-HMAC-SHA256" + " Credential=" + os.Getenv("AWS_ACCESS_KEY_ID") + "/" + credentialScope + " SignedHeaders=" + signedHeaders + " Signature=" + signature
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))

	if err != nil {
		fmt.Println("Request failed: ", err)
		return products
	}

	req.Header.Set("content-encoding", contentEncoding)
	req.Header.Set("content-type", contentType)
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

	json.NewDecoder(resp.Body).Decode(&products)
	return products
}
