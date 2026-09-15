package server

import (
	"log"
	"os"

	"sebosun/acrevus-go/internal/database"
)

type Envs struct {
	PORT   string
	DBURL string
	JwtSecret string
}

type APIConfig struct {
	DB *database.Queries
	ENVS Envs
}

func ReadEnvs() Envs {
	envs := Envs{
		PORT:   os.Getenv("PORT"),
		DBURL: os.Getenv("DATABASE_URL"),
		JwtSecret: os.Getenv("JWT_SECRET"),
	}

	if envs.PORT == "" {
		log.Fatal("PORT environment variable is not set")
	}

	if envs.DBURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	if envs.JwtSecret == "" {
		log.Fatal("JWT Secret nvironment variable is not set")
	}

	return envs
}
