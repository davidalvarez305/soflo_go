package types

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
