package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/davidalvarez305/soflo_go/server/types"
)

func FilterCommercialKeywords(results types.GoogleKeywordResults, seedKeyword string) []string {
	var data []string
	r := regexp.MustCompile("(used|cheap|deals|deal|sale|buy|online|on sale|discount|for sale|near me|best|for|[0-9]+)")

	for i := 0; i < len(results.Results); i++ {
		cleanKeyword := strings.TrimSpace(r.ReplaceAllString(results.Results[i].Text, ""))
		fmt.Println(cleanKeyword)

		compIndex, errOne := strconv.Atoi(results.Results[i].KeywordIdeaMetrics.CompetitionIndex)
		if errOne != nil {
			return data
		}

		searchVol, errTwo := strconv.Atoi(results.Results[i].KeywordIdeaMetrics.AvgMonthlySearches)
		if errTwo != nil {
			return data
		}

		conditionOne := compIndex == 100
		conditionTwo := searchVol > 100
		conditionThree := len(strings.Split(strings.TrimSpace(cleanKeyword), seedKeyword)[0]) >= 2

		if conditionOne && conditionTwo && conditionThree {
			data = append(data, cleanKeyword)
		}
	}

	return data
}

func GetGoogleCredentials() (types.GoogleConfigData, error) {
	data := types.GoogleConfigData{}

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

func ParseGoogleSERP(r io.Reader) ([]string, error) {
	var keywords []string

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		fmt.Println("Error trying to parse document.")
		return keywords, err
	}

	doc.Find(".ULSxyf").Each(func(i int, s *goquery.Selection) {
		el := s.Find("a").Text()
		fmt.Printf("%+v\n", strings.Split(el, "?"))
	})
	return keywords, nil
}
