package main

import (
	"fmt"
	"os"
	"ventas-app/config"
	"ventas-app/database"
	"ventas-app/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Este valor ya no es necesario si usás GetDB por request
	env := os.Getenv("APP_ENV")

	// Si es CI, forzamos entorno "ci"
	if os.Getenv("CI") == "true" {
		env = "ci"
	}

	if env == "" {
		env = "qa" // fallback por defecto (Render QA)
	}

	fmt.Println("Iniciando backend en entorno:", env)

	config.LoadEnv(env) // carga env.<app_env>

	database.Connect() // <- conecta según env actual

	r := gin.Default()
	fmt.Println("Conexión establecida para QA y PROD", env)

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5174"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Env"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	routes.Setup(r)

	r.Run(":8080")
}
