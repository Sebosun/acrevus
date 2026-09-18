package server

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"sebosun/acrevus-go/analyzer"
	"sebosun/acrevus-go/helpers"
	"sebosun/acrevus-go/internal/database"

	"github.com/gin-gonic/gin"
)

type ArticleFetch struct {
	URL string `json:"url" binding:"omitempty"`
}

type ArticleResponse struct {
	ID        int64      `json:"id"`
	Title     string     `json:"title"`
	HTML      string     `json:"html"`
	Author    string     `json:"author"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func (config *APIConfig) FetchArticle(c *gin.Context) {
	var userForm ArticleFetch

	err := c.ShouldBindJSON(&userForm)
	userID := c.MustGet(KeysUserID).(int64)
	params := database.LinkArticleToUserParams{
		UserID:    int64(userID),
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No payload or invalid payload"})
		return
	}

	cleanedURL := userForm.URL
	article, err := config.DB.GetArticleByURL(c.Request.Context(), cleanedURL)

	// TODO: Finish binding this with users already existing relationship with resource
	// Happy path we already have it
	if err == nil {
		params.ArticleID = article.ID
		_, err := config.DB.LinkArticleToUser(c.Request.Context(), params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		}

		response := NewArticleResponse(article)
		c.JSON(http.StatusOK, gin.H{
			"article": response,
		})
		return
	}

	if !errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}

	// TODO: I will have to strip the url out of any bs so we can compare them directly
	art, err := analyzer.Run(userForm.URL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to parse the article..."})
		return
	}

	articleDBParams := database.CreateArticleParams{
		Author: helpers.NewNullString(art.Author),
		Url:    cleanedURL,
		Html:   art.RawHTML,
		Title:  helpers.NewNullString(art.Title),
	}

	saved, err := config.DB.CreateArticle(c.Request.Context(), articleDBParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong..."})
		return
	}

	params.ArticleID = saved.ID
	_, err = config.DB.LinkArticleToUser(c.Request.Context(), params)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
	}

	response := NewArticleResponse(saved)

	c.JSON(http.StatusOK, gin.H{
		"article": response,
	})
}

func NewArticleResponse(article database.Article) ArticleResponse {
	response := ArticleResponse{
		ID:        article.ID,
		HTML:      article.Html,
		CreatedAt: article.CreatedAt,
		UpdatedAt: article.UpdatedAt,
	}

	if article.Author.Valid {
		response.Author = article.Author.String
	}
	if article.Author.Valid {
		response.Title = article.Title.String
	}

	if article.DeletedAt.Valid {
		response.DeletedAt = &article.DeletedAt.Time
	}

	return response
}
