package actions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/davidalvarez305/soflo_go/server/types"
)

func PullDynamicContent() []types.DynamicContent {
	var content []types.DynamicContent
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
		fmt.Println("Error while querying the dynamic content endpoint.", err)
		return content
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&content)
	return content
}

func PullContentDictionary() []types.Dictionary {
	var content []types.Dictionary
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
		fmt.Println("Error while querying the dictionary endpoint.", err)
		return content
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&content)
	return content
}
