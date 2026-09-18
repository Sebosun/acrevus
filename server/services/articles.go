package services

import (
	"sebosun/acrevus-go/analyzer"
	"sebosun/acrevus-go/internal/database"
)

type ArticleService struct {
	DB *database.Queries
}

func (service *ArticleService) ArticleFetcherKurwa(url string) (analyzer.MainArticle, error) {
	art, err := analyzer.Run(url)
	if err != nil {
		return analyzer.MainArticle{}, err
	}

	return art, nil
}
