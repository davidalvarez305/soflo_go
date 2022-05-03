package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/davidalvarez305/soflo_go/server/types"
)

func ParseHtml(r io.Reader, keyword string) ([]types.AmazonSearchResultsPage, error) {
	var products []types.AmazonSearchResultsPage

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		fmt.Println("Error trying to parse document.")
		return products, err
	}

	doc.Find(".sg-col-inner").Each(func(i int, s *goquery.Selection) {
		var product types.AmazonSearchResultsPage

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

func FormatShortDate(dt time.Time) string {
	return dt.UTC().Format("20060102")
}

func FormatDate(dt time.Time) string {
	return dt.UTC().Format("20060102T150405Z")
}

func MakeHash(hash hash.Hash, b []byte) []byte {
	hash.Reset()
	hash.Write(b)
	return hash.Sum(nil)
}

func BuildCanonicalString(method, uri, query, signedHeaders, canonicalHeaders, payloadHash string) string {
	return strings.Join([]string{
		method,
		uri,
		query,
		canonicalHeaders + "\n",
		signedHeaders,
		payloadHash,
	}, "\n")
}

func BuildStringToSign(date, credentialScope, canonicalRequest string) string {
	return strings.Join([]string{
		"AWS4-HMAC-SHA256",
		date,
		credentialScope,
		hex.EncodeToString(MakeHash(sha256.New(), []byte(canonicalRequest))),
	}, "\n")
}

func BuildCanonicalHeaders(contentType, contentEncoding, host, xAmazonDate, amazonTarget string) string {
	return strings.Join([]string{"content-encoding:" + contentEncoding, "content-type:" + contentType, "host:" + host, "x-amz-date:" + xAmazonDate, "x-amz-target:" + amazonTarget}, "\n")
}

func BuildSignature(strToSign string, sig []byte) (string, error) {
	return hex.EncodeToString(HMACSHA256(sig, []byte(strToSign))), nil
}

func HMACSHA256(key []byte, data []byte) []byte {
	hash := hmac.New(sha256.New, key)
	hash.Write(data)
	return hash.Sum(nil)
}
