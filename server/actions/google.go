package actions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Data struct {
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

func GetGoogleAuthToken() (Data, error) {
	data := Data{}

	path := "/home/david/soflo_go/server/google.json"

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

	fmt.Println(data)
	return data, nil
}
