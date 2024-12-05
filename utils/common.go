package utils

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load("dev.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}
