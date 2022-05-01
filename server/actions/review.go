package actions

import (
	"fmt"

	"github.com/davidalvarez305/soflo_go/server/types"
)

func GetProducts(keyword string) []types.AmazonSearchResultsPage {
	var products []types.AmazonSearchResultsPage

	q := types.GoogleQuery{
		Pagesize: 50,
		KeywordSeed: types.KeywordSeed{
			Keywords: [1]string{keyword},
		},
	}

	KW_ARR := QueryGoogle(q)

	seedKeywords := GetSeedKeywords(KW_ARR)
	commercialKeywords := GetCommercialKeywords(seedKeywords)

	for i := 0; i < len(commercialKeywords); i++ {
		data := ScrapeSearchResultsPage(commercialKeywords[i])
		if len(data) == 0 {
			fmt.Println("Keyword: " + commercialKeywords[i] + "0")
		}
		if len(data) > 0 {
			products = append(products, data...)
		}
		total := fmt.Sprintf("Keyword #%v of %v - %s - Total Products = %v", i+1, len(commercialKeywords), commercialKeywords[i], len(data))
		fmt.Println(total)
	}

	productsTotal := fmt.Sprintf("Total Products = %v", len(products))
	fmt.Println(productsTotal)

	return products
}

func PersistProducts(products []types.AmazonSearchResultsPage) bool {

	return true
}
