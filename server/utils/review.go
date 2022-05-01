package utils

import (
	"regexp"
	"strings"

	"github.com/davidalvarez305/soflo_go/server/models"
	"github.com/davidalvarez305/soflo_go/server/types"
)

func CreateCategorySlug(str string) string {
	var final string
	r := regexp.MustCompile(`[a-z0-9]+`)
	res := r.FindAllString(str, -1)

	if len(res) > 0 {
		final = strings.Join(res, "-")
	}

	return final
}

func CreateReviewPostFields(input types.AmazonSearchResultsPage, dictionary []types.Dictionary, sentences []types.DynamicContent, categoryId int) models.ReviewPost {
	name := strings.Join(strings.Split(strings.ToLower(input.Name), " "), "-")
	slug := CreateCategorySlug(name)

	data := GenerateContentUtil(input.Name, dictionary, sentences)

	fields := models.ReviewPost{
		Title:                         data.ReviewPostTitle,
		CategoryID:                    categoryId,
		Slug:                          slug,
		Content:                       data.ReviewPostContent,
		Headline:                      data.ReviewPostHeadline,
		Intro:                         data.ReviewPostIntro,
		Description:                   data.ReviewPostDescription,
		ProductLabel:                  data.ReviewPostProductLabel,
		ProductName:                   input.Name,
		ProductDescription:            data.ReviewPostProductDescription,
		ProductAffiliateUrl:           input.Link,
		Faq_Answer_1:                  data.ReviewPostFaq_Answer_1,
		Faq_Answer_2:                  data.ReviewPostFaq_Answer_2,
		Faq_Answer_3:                  data.ReviewPostFaq_Answer_3,
		Faq_Question_1:                data.ReviewPostFaq_Question_1,
		Faq_Question_2:                data.ReviewPostFaq_Question_2,
		Faq_Question_3:                data.ReviewPostFaq_Question_3,
		HorizontalCardProductImageUrl: input.Image,
		HorizontalCardProductImageAlt: strings.ToLower(input.Name),
	}

	return fields
}
