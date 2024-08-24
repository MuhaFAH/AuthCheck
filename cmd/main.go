package main

import (
	"github.com/MuhaFAH/AuthCheck/handlers"
	"github.com/MuhaFAH/AuthCheck/pkg/storage"
	"github.com/MuhaFAH/AuthCheck/routes"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error when loading .env file: %s", err.Error())
	}

	db, err := storage.NewDB(storage.Config{
		Host:     os.Getenv("DATABASE_HOST"),
		Port:     os.Getenv("DATABASE_PORT"),
		Username: os.Getenv("DATABASE_USER"),
		Password: os.Getenv("DATABASE_PASSWORD"),
		Database: os.Getenv("DATABASE_DB"),
		SSLMode:  os.Getenv("DATABASE_SSLMODE"),
	})
	if err != nil {
		log.Fatalf("error when connecting to database: %s", err.Error())
	}

	app := handlers.App{DB: db}
	router := routes.NewRouter(&app)

	if err := router.Run(os.Getenv("SERVER_PORT")); err != nil {
		log.Fatalf("error when running server: %s", err.Error())
	}
}
