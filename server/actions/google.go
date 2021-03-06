package actions

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/davidalvarez305/soflo_go/server/types"
	"github.com/davidalvarez305/soflo_go/server/utils"
)

func RequestGoogleAuthToken() (string, error) {
	config, err := utils.GetGoogleCredentials()
	if err != nil {
		fmt.Println("Error getting Google credentials")
		return "", err
	}

	client := &http.Client{}

	url := config.AuthURI
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Println("Request failed: ", err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	q := req.URL.Query()
	q.Add("access_type", "offline")
	q.Add("approval_prompt", "force")
	q.Add("scope", "https://www.googleapis.com/auth/adwords")
	q.Add("client_id", config.OAuthClientID)
	q.Add("redirect_uri", config.RedirectURI[0])
	q.Add("response_type", "code")
	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error while getting auth token", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "Request was not successful.", nil
	}

	var data http.Response

	e := json.NewDecoder(resp.Body).Decode(&data)
	if e != nil {
		return "", e
	}

	return "Some URL", nil
}

func GetGoogleAccessToken(code string) (string, error) {
	config, err := utils.GetGoogleCredentials()
	if err != nil {
		fmt.Println("Error getting Google credentials")
		return "", err
	}

	client := &http.Client{}

	url := config.AuthURI
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Println("Request failed: ", err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	q := req.URL.Query()
	q.Add("code", code)
	q.Add("client_id", config.OAuthClientID)
	q.Add("client_secret", config.OAuthClientSecret)
	q.Add("redirect_uri", config.RedirectURI[0])
	q.Add("scope", "https://www.googleapis.com/auth/adwords")
	q.Add("grant_type", "authorization_code")
	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error while getting auth token", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "Request was not successful.", nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error while reading response body: ", err)
		return "", err
	}

	data := string(body)

	return data, nil

}

func RefreshAuthToken() (string, error) {

	type TokenResponse struct {
		Access_Token string `json:"access_token"`
		Expires_In   string `json:"expires_in"`
		Scope        string `json:"scope"`
		Token_Type   string `json:"token_type"`
	}

	config, err := utils.GetGoogleCredentials()
	if err != nil {
		fmt.Println("Error getting Google credentials")
		return "", err
	}

	refreshToken := os.Getenv("REFRESH_TOKEN")
	client := &http.Client{}

	url := config.TokenURI
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Println("Request failed: ", err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	q := req.URL.Query()
	q.Add("client_id", config.OAuthClientID)
	q.Add("client_secret", config.OAuthClientSecret)
	q.Add("refresh_token", refreshToken)
	q.Add("grant_type", "refresh_token")
	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error while getting auth token", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "Request was not successful.", nil
	}

	var data TokenResponse

	json.NewDecoder(resp.Body).Decode(&data)

	return data.Access_Token, nil
}

func GetSeedKeywords(results types.GoogleKeywordResults) []string {
	var data []string

	for i := 0; i < len(results.Results); i++ {
		compIndex, errOne := strconv.Atoi(results.Results[i].KeywordIdeaMetrics.CompetitionIndex)
		if errOne != nil {
			return data
		}

		searchVol, errTwo := strconv.Atoi(results.Results[i].KeywordIdeaMetrics.AvgMonthlySearches)
		if errTwo != nil {
			return data
		}

		keywordLength := len(strings.Split(results.Results[i].Text, " "))

		conditionOne := compIndex == 100
		conditionTwo := searchVol > 10000
		conditionThree := keywordLength >= 2 && keywordLength <= 4

		if conditionOne && conditionTwo && conditionThree {
			data = append(data, results.Results[i].Text)
		}
	}

	fmt.Println("Seed Keywords: ", len(data))

	return data
}

func QueryGoogle(query types.GoogleQuery) types.GoogleKeywordResults {
	time.Sleep(1 * time.Second)
	var results types.GoogleKeywordResults

	authToken, err := RefreshAuthToken()

	if err != nil {
		fmt.Printf("Error refreshing token.")
		return results
	}

	googleCustomerID := os.Getenv("GOOGLE_CUSTOMER_ID")
	googleUrl := fmt.Sprintf("https://googleads.googleapis.com/v10/customers/%s:generateKeywordIdeas", googleCustomerID)
	developerToken := os.Getenv("GOOGLE_DEVELOPER_TOKEN")
	authorizationHeader := fmt.Sprintf("Bearer %s", authToken)

	client := &http.Client{}

	out, e := json.Marshal(query)

	if e != nil {
		return results
	}

	req, err := http.NewRequest("POST", googleUrl, bytes.NewBuffer(out))
	if err != nil {
		fmt.Println("Request failed: ", err)
		return results
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("developer-token", developerToken)
	req.Header.Set("Authorization", authorizationHeader)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error while querying Google", err)
		return results
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&results)

	return results
}

func GetCommercialKeywords(seedKeywords []string) []string {
	var keywords []string
	for i := 0; i < len(seedKeywords); i++ {

		q := types.GoogleQuery{
			Pagesize: 1000,
			KeywordSeed: types.KeywordSeed{
				Keywords: [1]string{seedKeywords[i]},
			},
		}

		results := QueryGoogle(q)
		k := utils.FilterCommercialKeywords(results, seedKeywords[i])
		keywords = append(keywords, k...)
	}

	fmt.Println("Commercial Keywords: ", len(keywords))

	return keywords
}

func CrawlGoogleSERP(keywords string) []string {
	var results []string

	str := strings.Join(strings.Split(keywords, " "), "+")

	serp := fmt.Sprintf("https://www.google.com/search?q=%s", str)

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
		fmt.Println("Error while fetching Google SERP", err)
		return results
	}
	defer resp.Body.Close()

	kws, err := utils.ParseGoogleSERP(resp.Body)
	results = kws

	if err != nil {
		fmt.Println("Error while parsing HTML.")
		return kws
	}

	return results
}
