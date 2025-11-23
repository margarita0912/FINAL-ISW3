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
		filename = ".env"
	}

	// Intentar cargar el archivo .env correspondiente.
	// Si no existe, NO abortamos (CI y PROD usan variables de entorno).
	err := godotenv.Load(filename)
	if err != nil {
		log.Printf("Archivo %s no encontrado o no se pudo cargar (no es un error en CI/PROD): %v", filename, err)
	}
}
