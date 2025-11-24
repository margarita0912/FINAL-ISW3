package database

import (
	"os"

	"github.com/gin-gonic/gin"
)

var GetDB func(c *gin.Context) DBHandler = func(c *gin.Context) DBHandler {

	// 1. Si estamos en CI, usamos la DB de CI directamente
	if os.Getenv("CI") == "true" {
		db := DBs["ci"]
		if db == nil {
			panic("Base de datos CI no inicializada")
		}
		return &GormDB{DB: db}
	}

	// 2. Caso contrario: seguir usando selector por header (tu lógica original)
	env := c.GetHeader("X-Env")
	if env != "prod" {
		env = "qa"
	}

	db := DBs[env]
	if db == nil {
		panic("Base de datos no inicializada para entorno: " + env)
	}

	return &GormDB{DB: db}
}
