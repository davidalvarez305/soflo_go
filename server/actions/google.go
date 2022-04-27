package actions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type GoogleConfigData struct {
	Type                    string   `json:"type"`
	ProjectID               string   `json:"project_id"`
	ProjectKeyId            string   `json:"private_key_id"`
	PrivateKey              string   `json:"private_key"`
	ClientEmail             string   `json:"client_email"`
	ClientID                string   `json:"client_id"`
	AuthURI                 string   `json:"auth_uri"`
	TokenURI                string   `json:"token_uri"`
	AuthProviderX509CertURL string   `json:"auth_provider_x509_cert_url"`
	ClientX509CertURL       string   `json:"client_x509_cert_url"`
	OAuthClientID           string   `json:"oauth_client_id"`
	OAuthClientSecret       string   `json:"oauth_client_secret"`
	RedirectURI             []string `json:"redirect_uris"`
	JavascriptOrigins       []string `json:"javascript_origins"`
}

type KeywordSeed struct {
	Keywords [1]string `json:"keywords"`
}

type GoogleQuery struct {
	Pagesize    int         `json:"pageSize"`
	KeywordSeed KeywordSeed `json:"keywordSeed"`
}

type MonthlySearchVolume struct {
	Month           string `json:"month"`
	Year            string `json:"year"`
	MonthlySearches string `json:"monthlySearches"`
}

type keywordIdeaMetrics struct {
	Competition            string                `json:"competition"`
	MonthlySearchVolume    []MonthlySearchVolume `json:"monthlySearchVolumes"`
	AvgMonthlySearches     string                `json:"avgMonthlySearches"`
	CompetitionIndex       string                `json:"competitionIndex"`
	LowTopOfPageBidMicros  string                `json:"lowTopOfPageBidMicros"`
	HighTopOfPageBidMicros string                `json:"highTopOfPageBidMicros"`
}

type GoogleResult struct {
	KeywordIdeaMetrics keywordIdeaMetrics `json:"keywordIdeaMetrics"`
	Text               string             `json:"text"`
}

type GoogleKeywordResults struct {
	Results []GoogleResult `json:"results"`
}

func GetGoogleCredentials() (GoogleConfigData, error) {
	data := GoogleConfigData{}

	path := os.Getenv("GOOGLE_JSON_PATH")

	file, err := os.Open(path)

	if err != nil {
		fmt.Println("Could not open Google JSON file.")
		return data, err
	}
	defer file.Close()

	jsonData, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Could not read Google JSON file.")
		return data, err
	}

	if err := json.Unmarshal(jsonData, &data); err != nil {
		fmt.Println("Error while trying to unmarshall JSON data.")
		return data, err
	}

	return data, nil
}

func RequestGoogleAuthToken() (string, error) {
	config, err := GetGoogleCredentials()
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
	config, err := GetGoogleCredentials()
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

	config, err := GetGoogleCredentials()
	if err != nil {
		fmt.Println("Error getting Google credentials")
		return "", err
	}

	refresh_token := os.Getenv("REFRESH_TOKEN")
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
	q.Add("refresh_token", refresh_token)
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

func filterSeedKeywords(results GoogleKeywordResults) []string {
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

		conditionOne := compIndex > 75
		conditionTwo := searchVol > 10000
		conditionThree := keywordLength >= 2 && keywordLength <= 3

		if conditionOne && conditionTwo && conditionThree {
			data = append(data, results.Results[i].Text)
		}
	}

	return data
}

func QueryGoogle(query GoogleQuery) []string {
	var results GoogleKeywordResults
	var keywords []string

	authToken, err := RefreshAuthToken()

	if err != nil {
		fmt.Printf("Error refreshing token.")
		return keywords
	}

	googleCustomerID := os.Getenv("GOOGLE_CUSTOMER_ID")
	googleUrl := fmt.Sprintf("https://googleads.googleapis.com/v10/customers/%s:generateKeywordIdeas", googleCustomerID)
	developerToken := os.Getenv("GOOGLE_DEVELOPER_TOKEN")
	authorizationHeader := fmt.Sprintf("Bearer %s", authToken)

	client := &http.Client{}

	out, e := json.Marshal(query)

	if e != nil {
		return keywords
	}

	req, err := http.NewRequest("POST", googleUrl, bytes.NewBuffer(out))
	if err != nil {
		fmt.Println("Request failed: ", err)
		return keywords
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("developer-token", developerToken)
	req.Header.Set("Authorization", authorizationHeader)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error while querying Google", err)
		return keywords
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&results)

	data := filterSeedKeywords(results)

	return data
}
