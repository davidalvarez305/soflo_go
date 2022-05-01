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

type ProcessedContent struct {
	ReviewPostTitle              []string `json:"ReviewPostTitle"`
	ReviewPostContent            []string `json:"ReviewPostContent"`
	ReviewPostHeadline           []string `json:"ReviewPostHeadline"`
	ReviewPostIntro              []string `json:"ReviewPostIntro"`
	ReviewPostDescription        []string `json:"ReviewPostDescription"`
	ReviewPostProductLabel       []string `json:"ReviewPostProductLabel"`
	ReviewPostProductDescription []string `json:"ReviewPostProductDescription"`
	ReviewPostFaq_Answer_1       []string `json:"ReviewPostFaq_Answer_1"`
	ReviewPostFaq_Answer_2       []string `json:"ReviewPostFaq_Answer_2"`
	ReviewPostFaq_Answer_3       []string `json:"ReviewPostFaq_Answer_3"`
	ReviewPostFaq_Question_1     []string `json:"ReviewPostFaq_Question_1"`
	ReviewPostFaq_Question_2     []string `json:"ReviewPostFaq_Question_2"`
	ReviewPostFaq_Question_3     []string `json:"ReviewPostFaq_Question_3"`
}

func filterSentences(sentence []actions.DynamicContent, paragraph string) []string {
	var s []string
	for i := 0; i < len(sentence); i++ {
		if sentence[i].Paragraph == paragraph {
			s = append(s, sentence[i].Content)
		}
	}
	return s
}

func GenerateContentUtil(dictionary []actions.Dictionary, sentences []actions.DynamicContent) []ProcessedContent {
	var dict []ProcessedDictionary
	var content []ProcessedContent

	for i := 0; i < len(dictionary); i++ {
		var d = ProcessedDictionary{
			Word: dictionary[i].Word,
			Tag:  strings.Split(dictionary[i].Content, "///"),
		}
		dict = append(dict, d)
	}
	fmt.Printf("%v", len(dict))

	for i := 0; i < len(sentences); i++ {
		a := filterSentences(sentences, "ReviewPostTitle")
		b := filterSentences(sentences, "ReviewPostContent")
		c := filterSentences(sentences, "ReviewPostHeadline")
		d := filterSentences(sentences, "ReviewPostIntro")
		e := filterSentences(sentences, "ReviewPostDescription")
		f := filterSentences(sentences, "ReviewPostProductLabel")
		g := filterSentences(sentences, "ReviewPostProductDescription")
		h := filterSentences(sentences, "ReviewPostFaq_Answer_1")
		j := filterSentences(sentences, "ReviewPostFaq_Answer_2")
		k := filterSentences(sentences, "ReviewPostFaq_Answer_3")
		l := filterSentences(sentences, "ReviewPostFaq_Question_1")
		m := filterSentences(sentences, "ReviewPostFaq_Question_2")
		n := filterSentences(sentences, "ReviewPostFaq_Question_3")
		var final = ProcessedContent{
			ReviewPostTitle:              a,
			ReviewPostContent:            b,
			ReviewPostHeadline:           c,
			ReviewPostIntro:              d,
			ReviewPostDescription:        e,
			ReviewPostProductLabel:       f,
			ReviewPostProductDescription: g,
			ReviewPostFaq_Answer_1:       h,
			ReviewPostFaq_Answer_2:       j,
			ReviewPostFaq_Answer_3:       k,
			ReviewPostFaq_Question_1:     l,
			ReviewPostFaq_Question_2:     m,
			ReviewPostFaq_Question_3:     n,
		}
		content = append(content, final)
	}
	return content
}
