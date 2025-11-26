package main

import (
	"fmt"
	"os"
	"path/filepath"
	"ventas-app/config"
	"ventas-app/database"
	"ventas-app/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	// Detectar entorno
	env := os.Getenv("APP_ENV")

	// Si es CI → usar entorno CI
	if os.Getenv("CI") == "true" {
		env = "ci"
	}

	// Si no viene nada (Render), usar QA como default
	if env == "" {
		env = "qa"
	}

	fmt.Println("Iniciando backend en entorno:", env)

	// Cargar variables desde env.<APP_ENV>
	config.LoadEnv(env)

	// Conectar BD según entorno
	database.Connect()

	r := gin.Default()

	fmt.Println("Conexión establecida para entorno:", env)

	// ===========================
	//        🔥 CORS FINAL
	// ===========================
	allowedOrigins := []string{
		// Localhost para Vite y Cypress
		"http://localhost:5173",
		"http://localhost:5174",

		// QA
		"https://frontqa-t0a9.onrender.com",

		// PROD (OJO: ESTA es la correcta según Render)
		"https://frontqa.onrender.com",
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Registrar rutas de API
	routes.Setup(r)

	// Servir frontend estático desde ./dist si existe
	// Esto permite que un único servicio de Render atienda frontend y backend.
	distPath := filepath.Join(".", "dist")
	r.Static("/assets", filepath.Join(distPath, "assets"))
	// Rutas no manejadas por la API → devolver index.html (SPA fallback)
	r.NoRoute(func(c *gin.Context) {
		c.File(filepath.Join(distPath, "index.html"))
	})

	// Puerto: usar PORT de Render si está definido, si no 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
