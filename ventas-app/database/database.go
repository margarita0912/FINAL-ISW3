package database

import (
	"fmt"
	"log"
	"ventas-app/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DBs = make(map[string]*gorm.DB)

// ConnectAll carga las conexiones a QA y PROD desde los archivos .env.qa y .env.prod
/*func ConnectAll() {
	// Cargar .env.qa
	envQA, err := godotenv.Read(".env.qa")
	if err != nil {
		log.Fatal("❌ Error cargando .env.qa:", err)
	}
	DBs["qa"] = connectFromMap(envQA)

	// Cargar .env.prod
	envProd, err := godotenv.Read(".env.prod")
	if err != nil {
		log.Fatal("❌ Error cargando .env.prod:", err)
	}
	DBs["prod"] = connectFromMap(envProd)

	fmt.Println("✅ Conexiones inicializadas:")
	for key, db := range DBs {
		if db != nil {
			fmt.Println("✔️", key, "conectado")
		} else {
			fmt.Println("❌", key, "no conectado")
		}
	}
} */

func connectFromMap(env map[string]string) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		env["DB_USER"],
		env["DB_PASS"],
		env["DB_HOST"],
		env["DB_PORT"],
		env["DB_NAME"],
	)

	fmt.Printf("🔗 Conectando a %s (%s:%s)\n", env["DB_NAME"], env["DB_HOST"], env["DB_PORT"])

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Error al conectar a %s: %v", env["DB_NAME"], err)
	}

	db.AutoMigrate(&models.Usuario{}, &models.Producto{}, &models.Compra{}, &models.Venta{})
	return db
}
