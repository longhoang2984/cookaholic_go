package infrastructure

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}
