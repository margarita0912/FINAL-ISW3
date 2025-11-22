package controllers

import (
	"net/http"
	"ventas-app/models"
	"ventas-app/utils"

	"ventas-app/database"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginInput struct {
	Nombre string `json:"nombre"`
	Clave  string `json:"clave"`
}

func Login(c *gin.Context) {
	db := database.GetDB(c) // ← función centralizada en database/db_selector.go

	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	var user models.Usuario
	if err := db.Where("nombre = ?", input.Nombre).First(&user); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no encontrado"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Clave), []byte(input.Clave)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Clave incorrecta"})
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Rol)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo generar el token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"rol":   user.Rol, // ← esto depende de tu modelo
	})

}
