package server

import (
	"net/http"

	"sebosun/acrevus-go/analyzer"

	"github.com/gin-gonic/gin"
)

type ArticleFetch struct {
	URL string `json:"url" binding:"omitempty"`
}

func (config *APIConfig) FetchArticle(c *gin.Context) {
	var userForm ArticleFetch

	err := c.ShouldBindJSON(&userForm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid url"})
		return
	}

	art, err := analyzer.Run(userForm.URL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to parse the article..."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"html":   art.RawHTML,
		"title":  art.Title,
		"author": art.Author,
	})
}
