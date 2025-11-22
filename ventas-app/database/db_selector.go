package database

import (
	"github.com/gin-gonic/gin"
)

// GetDB es una variable que devuelve un database.DBHandler
// Se declara como variable para que los tests puedan sobrescribirla con mocks.
var GetDB func(c *gin.Context) DBHandler = func(c *gin.Context) DBHandler {
	env := c.GetHeader("X-Env")
	if env != "prod" {
		env = "qa"
	}
	db := DBs[env]
	if db == nil {
		panic("Base de datos no inicializada para entorno: " + env)
	}
	// Envolvemos *gorm.DB en GormDB para cumplir DBHandler
	return &GormDB{DB: db}
}
