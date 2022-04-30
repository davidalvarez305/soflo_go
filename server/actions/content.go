package actions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type DynamicContent struct {
	Content   string `json:"content"`
	Template  string `json:"template"`
	Paragraph string `json:"paragraph"`
}
type Dictionary struct {
	Word    string `json:"word"`
	Tag     string `json:"tag"`
	Content string `json:"content"`
}

func PullDynamicContent() []DynamicContent {
	var content []DynamicContent
	contentApi := os.Getenv("DYNAMIC_CONTENT_API") + "/api/get-dynamic-content/?template=ReviewPost"

	client := &http.Client{}
	req, err := http.NewRequest("GET", contentApi, nil)

	if err != nil {
		fmt.Println("Request failed: ", err)
		return content
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", os.Getenv("GENERATOR_TOKEN")))

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error while querying Google", err)
		return content
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&content)
	return content
}

func PullContentDictionary() []Dictionary {
	var content []Dictionary
	contentApi := os.Getenv("DYNAMIC_CONTENT_API") + "/api/get-dictionary/"

	client := &http.Client{}
	req, err := http.NewRequest("GET", contentApi, nil)

	if err != nil {
		fmt.Println("Request failed: ", err)
		return content
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", os.Getenv("GENERATOR_TOKEN")))

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error while querying Google", err)
		return content
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&content)
	return content
}
