package types

type PAAAPI5Response struct {
	SearchResult SearchResult `json:"SearchResult"`
}

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
