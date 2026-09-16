// Package server handles the server, duh!
package server

import (
	"database/sql"
	"log"

	"sebosun/acrevus-go/internal/database"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

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
	}

	apiRoute := router.Group("/api/v1")

	apiRoute.POST("/login", apiConfig.UserLogin)
	apiRoute.POST("/register", apiConfig.UserRegister)

	apiRoute.GET("/users", apiConfig.UserList)
	apiRoute.GET("/users/:id", apiConfig.UserGet)
	apiRoute.PATCH("/users/:id", apiConfig.UserUpdate)
	apiRoute.DELETE("/users/:id", apiConfig.UserDelete)

	authorizedRoutes := router.Group("")
	authorizedRoutes.Use(apiConfig.AuthRequired())

	authorizedRoutes.GET("/me", apiConfig.FetchMe)

	apiRoute.POST("/article", apiConfig.FetchArticle)

	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	router.Run()
}
