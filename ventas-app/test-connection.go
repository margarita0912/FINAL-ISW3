package main

import (
	"fmt"
	"log"
	"ventas-app/database"

	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: No se pudo cargar el archivo .env: %v", err)
	}

	fmt.Println("Intentando conectar a la base de datos Aiven...")

	// Intentar conectar
	db := database.Connect()
	if db != nil {
		fmt.Println("✅ Conexión exitosa a Aiven MySQL!")

		// Verificar que podemos hacer una consulta simple
		var version string
		err := db.Raw("SELECT VERSION()").Scan(&version).Error
		if err != nil {
			log.Printf("Error ejecutando consulta de prueba: %v", err)
		} else {
			fmt.Printf("🎉 Versión de MySQL: %s\n", version)
		}
	} else {
		fmt.Println("❌ Error: no se pudo establecer la conexión")
	}
}
