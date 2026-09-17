// Package server handles the server, duh!
package server

import (
	"database/sql"
	"log"

	"sebosun/acrevus-go/internal/database"
	"sebosun/acrevus-go/server/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// TODO: Some service for sending back errors

func StartServer() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
		return
	}

	// Create a Gin router with default middleware (logger and recovery)
	router := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AddAllowHeaders("Authorization")

	router.Use(cors.New(corsConfig))

	envs := ReadEnvs()

	db, err := sql.Open("postgres", envs.DBURL)
	if err != nil {
		log.Fatal(err)
		return
	}

	defer db.Close()

	dbQueries := database.New(db)

	apiConfig := APIConfig{
		DB:   dbQueries,
		ENVS: envs,
		services: Services{
			jwt: services.JwtService{
				Secret: envs.JwtSecret,
			},
			articles: services.ArticleService{},
		},
	}

	apiRoute := router.Group("/api/v1")
	authorizedRoutes := apiRoute.Group("")

	authorizedRoutes.Use(apiConfig.AuthRequired())
	apiConfig.RegisterAuthorizedRoutes(authorizedRoutes)
	apiConfig.RegisterPublicRoutes(apiRoute)

	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	router.Run()
}

func (config *APIConfig) RegisterPublicRoutes(router *gin.RouterGroup) {
	router.POST("/login", config.UserLogin)
	router.POST("/register", config.UserRegister)

	router.GET("/users", config.UserList)
	router.GET("/users/:id", config.UserGet)
	router.PATCH("/users/:id", config.UserUpdate)
	router.DELETE("/users/:id", config.UserDelete)
}

func (config *APIConfig) RegisterAuthorizedRoutes(router *gin.RouterGroup) {
	router.GET("/me", config.FetchMe)
	router.POST("/article", config.FetchArticle)
}
