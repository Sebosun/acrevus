package services

import "sebosun/acrevus-go/analyzer"

type ArticleService struct {
	Secret string
}

func (service *ArticleService) ArticleFetcherKurwa(url string) (analyzer.MainArticle, error) {
	art, err := analyzer.Run(url)
	if err != nil {
		return analyzer.MainArticle{}, err
	}

	return art, nil
}
