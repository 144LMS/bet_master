package initializers

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load("D:/vsCode/bet_master/.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}
