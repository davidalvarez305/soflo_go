package utils

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"

	"github.com/davidalvarez305/soflo_go/server/actions"
)

type ProcessedDictionary struct {
	Word    string   `json:"word"`
	Tag     string   `json:"tag"`
	Content []string `json:"content"`
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

type FinalizedContent struct {
	ReviewPostTitle              string `json:"ReviewPostTitle"`
	ReviewPostContent            string `json:"ReviewPostContent"`
	ReviewPostHeadline           string `json:"ReviewPostHeadline"`
	ReviewPostIntro              string `json:"ReviewPostIntro"`
	ReviewPostDescription        string `json:"ReviewPostDescription"`
	ReviewPostProductLabel       string `json:"ReviewPostProductLabel"`
	ReviewPostProductDescription string `json:"ReviewPostProductDescription"`
	ReviewPostFaq_Answer_1       string `json:"ReviewPostFaq_Answer_1"`
	ReviewPostFaq_Answer_2       string `json:"ReviewPostFaq_Answer_2"`
	ReviewPostFaq_Answer_3       string `json:"ReviewPostFaq_Answer_3"`
	ReviewPostFaq_Question_1     string `json:"ReviewPostFaq_Question_1"`
	ReviewPostFaq_Question_2     string `json:"ReviewPostFaq_Question_2"`
	ReviewPostFaq_Question_3     string `json:"ReviewPostFaq_Question_3"`
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

func processSentence(productName, sentence string, dictionary []ProcessedDictionary) string {
	var s string

	r := regexp.MustCompile(`(\([#@]\w+:[A-Z]+)\)|(\([#@]\w+)\)`)

	m := r.FindAllString(sentence, -1)

	for i := 0; i < len(m); i++ {
		if i == 0 {
			s = sentence
		}
		switched := strings.Replace(s, m[i], spinnerFunction(productName, m[i], dictionary), -1)
		s = switched
	}
	return s
}

func switchWords(matchedWord string, dictionary []ProcessedDictionary) string {
	for i := 0; i < len(dictionary); i++ {
		if dictionary[i].Tag == matchedWord {
			matchedWord = dictionary[i].Content[rand.Intn(len(dictionary[i].Content))]
		}
	}
	return matchedWord

}

func spinnerFunction(productName, matchedWord string, dictionary []ProcessedDictionary) string {
	if matchedWord == "(@ProductName)" {
		matchedWord = productName
		return matchedWord
	} else {
		splitStr := strings.Split(matchedWord, ":")
		if len(splitStr) == 2 {
			s := splitStr[0] + ")"
			matchedWord = switchWords(s, dictionary)
			if splitStr[1] == "UU)" {
				matchedWord = strings.Title(matchedWord)
			}
			if splitStr[1] == "U)" {
				ss := strings.Split(matchedWord, "")
				ss[0] = strings.ToUpper(ss[0])
				matchedWord = strings.Join(ss, "")
			}
		} else {
			matchedWord = switchWords(matchedWord, dictionary)
		}
	}
	return matchedWord
}

func selectRandomSentences(productName string, sentences []ProcessedContent, dictionary []ProcessedDictionary) FinalizedContent {
	var content FinalizedContent
	for i := 0; i < len(sentences); i++ {
		content.ReviewPostTitle = processSentence(productName, sentences[i].ReviewPostTitle[rand.Intn(len(sentences[i].ReviewPostTitle))], dictionary)
		content.ReviewPostContent = processSentence(productName, sentences[i].ReviewPostContent[rand.Intn(len(sentences[i].ReviewPostContent))], dictionary)
		content.ReviewPostHeadline = processSentence(productName, sentences[i].ReviewPostHeadline[rand.Intn(len(sentences[i].ReviewPostHeadline))], dictionary)
		content.ReviewPostIntro = processSentence(productName, sentences[i].ReviewPostIntro[rand.Intn(len(sentences[i].ReviewPostIntro))], dictionary)
		content.ReviewPostDescription = processSentence(productName, sentences[i].ReviewPostDescription[rand.Intn(len(sentences[i].ReviewPostDescription))], dictionary)
		content.ReviewPostProductLabel = processSentence(productName, sentences[i].ReviewPostProductLabel[rand.Intn(len(sentences[i].ReviewPostProductLabel))], dictionary)
		content.ReviewPostProductDescription = processSentence(productName, sentences[i].ReviewPostProductDescription[rand.Intn(len(sentences[i].ReviewPostProductDescription))], dictionary)
		content.ReviewPostFaq_Answer_1 = processSentence(productName, sentences[i].ReviewPostFaq_Answer_1[rand.Intn(len(sentences[i].ReviewPostFaq_Answer_1))], dictionary)
		content.ReviewPostFaq_Answer_2 = processSentence(productName, sentences[i].ReviewPostFaq_Answer_2[rand.Intn(len(sentences[i].ReviewPostFaq_Answer_2))], dictionary)
		content.ReviewPostFaq_Answer_3 = processSentence(productName, sentences[i].ReviewPostFaq_Answer_3[rand.Intn(len(sentences[i].ReviewPostFaq_Answer_3))], dictionary)
		content.ReviewPostFaq_Question_1 = processSentence(productName, sentences[i].ReviewPostFaq_Question_1[rand.Intn(len(sentences[i].ReviewPostFaq_Question_1))], dictionary)
		content.ReviewPostFaq_Question_2 = processSentence(productName, sentences[i].ReviewPostFaq_Question_2[rand.Intn(len(sentences[i].ReviewPostFaq_Question_2))], dictionary)
		content.ReviewPostFaq_Question_3 = processSentence(productName, sentences[i].ReviewPostFaq_Question_3[rand.Intn(len(sentences[i].ReviewPostFaq_Question_3))], dictionary)
	}
	return content
}

func GenerateContentUtil(productName string, dictionary []actions.Dictionary, sentences []actions.DynamicContent) FinalizedContent {
	var dict []ProcessedDictionary
	var content []ProcessedContent
	var finalContent FinalizedContent

	for i := 0; i < len(dictionary); i++ {
		var d = ProcessedDictionary{
			Word:    dictionary[i].Word,
			Tag:     dictionary[i].Tag,
			Content: strings.Split(dictionary[i].Content, "///"),
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
	finalContent = selectRandomSentences(productName, content, dict)
	return finalContent
}
