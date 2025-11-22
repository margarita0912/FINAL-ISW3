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
	config.LoadEnv("qa") // Carga inicial mínima para variables comunes

	database.ConnectAll() // ← Carga ambas bases: qa y prod

	r := gin.Default()
	fmt.Println("Conexión establecida para QA y PROD", env)

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Env"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	routes.Setup(r)

	r.Run(":8080")
}
