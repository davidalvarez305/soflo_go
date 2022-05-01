package utils

import (
	"fmt"
	"strings"

	"github.com/davidalvarez305/soflo_go/server/actions"
)

type ProcessedDictionary struct {
	Word string   `json:"word"`
	Tag  []string `json:"tag"`
}

func GenerateContentUtil(dictionary []actions.Dictionary) []ProcessedDictionary {
	var dict []ProcessedDictionary

	for i := 0; i < len(dictionary); i++ {
		dict[i].Word = dictionary[i].Word
		dict[i].Tag = strings.Split(dictionary[i].Content, "///")
	}

	fmt.Println(dict)
	return dict
}
