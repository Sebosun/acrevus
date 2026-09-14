package server

import (
	"log"
	"os"

	"sebosun/acrevus-go/internal/database"
)

type Envs struct {
	PORT   string
	DB_URL string
	JWT_SECRET string
}

type ApiConfig struct {
	DB *database.Queries
	ENVS Envs
}

func ReadEnvs() Envs {
	envs := Envs{
		PORT:   os.Getenv("PORT"),
		DB_URL: os.Getenv("DATABASE_URL"),
		JWT_SECRET: os.Getenv("SECRET"),
	}

	if envs.PORT == "" {
		log.Fatal("PORT environment variable is not set")
	}

	if envs.DB_URL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	if envs.JWT_SECRET == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	return envs
}
