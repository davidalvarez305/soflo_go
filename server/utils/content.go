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

func processSentence(sentence string, dictionary []ProcessedDictionary) string {
	var s string

	r := regexp.MustCompile(`(\([#@]\w+:[A-Z]+)\)|(\([#@]\w+)\)`)

	m := r.FindAllString(sentence, 3)
	fmt.Printf("m: %s", m)

	for i := 0; i < len(m); i++ {
		switched := strings.Replace(sentence, m[i], spinnerFunction(m[i], dictionary), 1)
		s = switched
	}
	fmt.Printf("Finalized sentence: %s", s)
	return s
}

func switchWords(matchedWord string, dictionary []ProcessedDictionary) string {
	var str string

	for i := 0; i < len(dictionary); i++ {
		if dictionary[i].Tag == matchedWord {
			str = dictionary[i].Content[rand.Intn(len(dictionary[i].Content))]
		}
	}
	return str

}

func spinnerFunction(matchedWord string, dictionary []ProcessedDictionary) string {
	var str string

	splitStr := strings.Split(matchedWord, ":")
	if len(splitStr) == 2 {
		s := splitStr[0] + ")"
		str = switchWords(s, dictionary)
		if splitStr[1] == "UU)" {
			str = strings.Title(str)
		}
		if splitStr[1] == "U)" {
			ss := strings.Split(str, "")
			ss[0] = strings.ToUpper(ss[0])
			str = strings.Join(ss, "")
		}
	} else {
		str = switchWords(matchedWord, dictionary)
	}
	return str
}

func selectRandomSentences(sentences []ProcessedContent, dictionary []ProcessedDictionary) FinalizedContent {
	var content FinalizedContent
	for i := 0; i < len(sentences); i++ {
		content.ReviewPostTitle = sentences[i].ReviewPostTitle[rand.Intn(len(sentences[i].ReviewPostTitle))]
		content.ReviewPostContent = sentences[i].ReviewPostContent[rand.Intn(len(sentences[i].ReviewPostContent))]
		content.ReviewPostHeadline = sentences[i].ReviewPostHeadline[rand.Intn(len(sentences[i].ReviewPostHeadline))]
		content.ReviewPostIntro = sentences[i].ReviewPostIntro[rand.Intn(len(sentences[i].ReviewPostIntro))]
		content.ReviewPostDescription = sentences[i].ReviewPostDescription[rand.Intn(len(sentences[i].ReviewPostDescription))]
		content.ReviewPostProductLabel = sentences[i].ReviewPostProductLabel[rand.Intn(len(sentences[i].ReviewPostProductLabel))]
		content.ReviewPostProductDescription = sentences[i].ReviewPostProductDescription[rand.Intn(len(sentences[i].ReviewPostProductDescription))]
		content.ReviewPostFaq_Answer_1 = sentences[i].ReviewPostFaq_Answer_1[rand.Intn(len(sentences[i].ReviewPostFaq_Answer_1))]
		content.ReviewPostFaq_Answer_2 = sentences[i].ReviewPostFaq_Answer_2[rand.Intn(len(sentences[i].ReviewPostFaq_Answer_2))]
		content.ReviewPostFaq_Answer_3 = sentences[i].ReviewPostFaq_Answer_3[rand.Intn(len(sentences[i].ReviewPostFaq_Answer_3))]
		content.ReviewPostFaq_Question_1 = sentences[i].ReviewPostFaq_Question_1[rand.Intn(len(sentences[i].ReviewPostFaq_Question_1))]
		content.ReviewPostFaq_Question_2 = sentences[i].ReviewPostFaq_Question_2[rand.Intn(len(sentences[i].ReviewPostFaq_Question_2))]
		content.ReviewPostFaq_Question_3 = sentences[i].ReviewPostFaq_Question_3[rand.Intn(len(sentences[i].ReviewPostFaq_Question_3))]
	}
	processSentence(content.ReviewPostTitle, dictionary)
	return content
}

func GenerateContentUtil(dictionary []actions.Dictionary, sentences []actions.DynamicContent) FinalizedContent {
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
	finalContent = selectRandomSentences(content, dict)
	return finalContent
}
