package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func EnvMongoURI() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error Loading .env file")
	}

	return os.Getenv("MONGOURI")
}

func GetPort() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error Loading .env file")
	}

	return os.Getenv("PORT")
}

func GetDB() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error Loading .env file")
	}

	return os.Getenv("DB")
}
