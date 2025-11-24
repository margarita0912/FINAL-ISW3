package database

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"ventas-app/models"

	"github.com/go-sql-driver/mysql"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() *gorm.DB {
	usuario := os.Getenv("DB_USER")
	clave := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	nombre := os.Getenv("DB_NAME")
	sslMode := os.Getenv("DB_SSL_MODE")

	// Configurar SSL para Aiven con skip-verify como fallback (solo QA/PROD)
	err := configurarSSLAiven()
	if err != nil {
		log.Printf("Warning: No se pudo configurar SSL con certificado CA: %v", err)
		err = configurarSSLAivenSimple()
		if err != nil {
			log.Printf("Warning: No se pudo configurar SSL simple: %v", err)
		}
	}

	var db *gorm.DB

	// CI/local → sin TLS
	if sslMode == "disable" {
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?parseTime=true",
			usuario, clave, host, port, nombre,
		)
		db, err = gorm.Open(gormMysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatal("Error al conectar con MySQL (sin TLS):", err)
		}
	} else {
		// QA/PROD (Aiven con TLS)
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?parseTime=true&tls=custom",
			usuario, clave, host, port, nombre,
		)

		db, err = gorm.Open(gormMysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Printf("Error con SSL personalizado, intentando SSL básico: %v", err)

			dsn = fmt.Sprintf(
				"%s:%s@tcp(%s:%s)/%s?parseTime=true&tls=skip-verify",
				usuario, clave, host, port, nombre,
			)

			db, err = gorm.Open(gormMysql.Open(dsn), &gorm.Config{})
			if err != nil {
				log.Fatal("Error al conectar con MySQL (Aiven):", err)
			}
		}
	}

	// Migraciones comunes
	db.AutoMigrate(
		&models.Usuario{},
		&models.Producto{},
		&models.Compra{},
		&models.Venta{},
	)

	// 🔥 REGISTRAR CONEXIÓN EN EL MAPA DBs
	// SI ESTAMOS EN CI → usar clave "ci"
	env := os.Getenv("APP_ENV")
	if os.Getenv("CI") == "true" {
		env = "ci"
	}

	DBs[env] = db // <-- esta línea hace que GetDB funcione perfecto en CI/QA/PROD

	DB = db
	return db
}

func configurarSSLAiven() error {
	caCertPath := "BaltimoreCyberTrustRoot.crt.pem"
	caCert, err := os.ReadFile(caCertPath)
	if err != nil {
		return fmt.Errorf("error leyendo certificado CA: %v", err)
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return fmt.Errorf("error agregando certificado CA al pool")
	}

	tlsConfig := &tls.Config{
		RootCAs:            caCertPool,
		InsecureSkipVerify: false,
	}

	err = mysql.RegisterTLSConfig("custom", tlsConfig)
	if err != nil {
		return fmt.Errorf("error registrando configuración TLS: %v", err)
	}

	return nil
}

func configurarSSLAivenSimple() error {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	err := mysql.RegisterTLSConfig("custom", tlsConfig)
	if err != nil {
		return fmt.Errorf("error registrando configuración TLS simple: %v", err)
	}

	return nil
}
