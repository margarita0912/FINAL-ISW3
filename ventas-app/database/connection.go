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

	// Configurar SSL para Aiven con skip-verify como fallback
	err := configurarSSLAiven()
	if err != nil {
		log.Printf("Warning: No se pudo configurar SSL con certificado CA: %v", err)
		// Configurar SSL con skip-verify como fallback para Aiven
		err = configurarSSLAivenSimple()
		if err != nil {
			log.Printf("Warning: No se pudo configurar SSL simple: %v", err)
		}
	}

	// Configuración DSN para Aiven MySQL con SSL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&tls=custom",
		usuario, clave, host, port, nombre)

	db, err := gorm.Open(gormMysql.Open(dsn), &gorm.Config{})
	if err != nil {
		// Fallback: intentar con SSL básico si falla con certificado personalizado
		log.Printf("Error con SSL personalizado, intentando SSL básico: %v", err)
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&tls=skip-verify",
			usuario, clave, host, port, nombre)
		db, err = gorm.Open(gormMysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatal("Error al conectar con MySQL (Aiven):", err)
		}
	}

	db.AutoMigrate(&models.Usuario{}, &models.Producto{}, &models.Compra{}, &models.Venta{})
	DB = db
	return db
}

func configurarSSLAiven() error {
	// Cargar certificado CA
	caCertPath := "BaltimoreCyberTrustRoot.crt.pem"
	caCert, err := os.ReadFile(caCertPath)
	if err != nil {
		return fmt.Errorf("error leyendo certificado CA: %v", err)
	}

	// Crear pool de certificados
	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return fmt.Errorf("error agregando certificado CA al pool")
	}

	// Configurar TLS
	tlsConfig := &tls.Config{
		RootCAs:            caCertPool,
		InsecureSkipVerify: false,
	}

	// Registrar configuración TLS personalizada
	err = mysql.RegisterTLSConfig("custom", tlsConfig)
	if err != nil {
		return fmt.Errorf("error registrando configuración TLS: %v", err)
	}

	return nil
}

func configurarSSLAivenSimple() error {
	// Configurar TLS con skip verify para Aiven (menos seguro pero funcional)
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	// Registrar configuración TLS personalizada
	err := mysql.RegisterTLSConfig("custom", tlsConfig)
	if err != nil {
		return fmt.Errorf("error registrando configuración TLS simple: %v", err)
	}

	return nil
}
