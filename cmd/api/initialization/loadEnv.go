package initialization

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv(filepaths ...string) {
	err := godotenv.Load(filepaths...)
	if err != nil {
		log.Printf("Error loading env file: %s", err.Error())
		os.Exit(1)
	}
}
