package actions

import (
	"fmt"

	"github.com/davidalvarez305/soflo_go/server/types"
	"github.com/davidalvarez305/soflo_go/server/utils"
)

func CreateReviewPosts(keyword string) ([]types.AmazonSearchResultsPage, error) {
	var products []types.AmazonSearchResultsPage
	dictionary := PullContentDictionary()
	sentences := PullDynamicContent()

	q := types.GoogleQuery{
		Pagesize: 1000,
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
			err := utils.InsertReviewPosts(data, dictionary, sentences)

			if err != nil {
				fmt.Printf("Error while trying to insert %s: %+v", commercialKeywords[i], err)
			}
			products = append(products, data...)
		}
		total := fmt.Sprintf("Keyword #%v of %v - %s - Total Products = %v", i+1, len(commercialKeywords), commercialKeywords[i], len(data))
		fmt.Println(total)
	}

	productsTotal := fmt.Sprintf("Total Products = %v", len(products))
	fmt.Println(productsTotal)

	return products, nil
}
