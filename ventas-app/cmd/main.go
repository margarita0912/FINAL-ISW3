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

	// ===========================
	// 🔍 DETECTAR ENTORNO
	// ===========================
	env := os.Getenv("APP_ENV")

	if os.Getenv("CI") == "true" {
		env = "ci"
	}

	if env == "" {
		env = "qa"
	}

	fmt.Println("Iniciando backend en entorno:", env)

	// ===========================
	// 🔧 CARGA DE CONFIG POR ENTORNO
	// ===========================
	config.LoadEnv(env)
	database.Connect()

	r := gin.Default()
	fmt.Println("Conexión establecida para entorno:", env)

	// ===========================
	// 🔥 CORS DINÁMICO POR ENTORNO
	// ===========================
	allowedOrigins := []string{
		"http://localhost:5173",
		"http://localhost:5174",
	}

	// Agregar los orígenes según entorno
	switch env {
	case "qa":
		allowedOrigins = append(allowedOrigins,
			"https://frontqa-t0a9.onrender.com",
		)
	case "prod":
		allowedOrigins = append(allowedOrigins,
			"https://frontprod-tn48.onrender.com", // <-- DOMINIO REAL DE PROD
		)
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// ===========================
	// 🔀 Rutas
	// ===========================
	routes.Setup(r)

	// ===========================
	// 🚀 Ejecutar servidor
	// ===========================
	r.Run(":8080")
}
