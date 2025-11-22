package config

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnv(env string) {
	var filename string
	switch env {
	case "prod":
		filename = ".env.prod"
	case "qa":
		filename = ".env.qa"
	default:
		filename = ".env.qa"
	}

	err := godotenv.Load(filename)
	if err != nil {
		log.Fatalf("Error cargando archivo %s: %v", filename, err)
	}
}
