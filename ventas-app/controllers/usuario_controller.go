package controllers

import (
	"net/http"
	"strings"
	"ventas-app/database"
	"ventas-app/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var rolesValidos = map[string]bool{
	"vendedor":  true,
	"comprador": true,
	"precio":    true,
}

type CrearUsuarioInput struct {
	Nombre string `json:"nombre" binding:"required"`
	Clave  string `json:"clave" binding:"required"`
	Rol    string `json:"rol" binding:"required"`
}

func CrearUsuario(c *gin.Context) {
	db := database.GetDB(c) // ← selecciona la base según el entorno

	var input CrearUsuarioInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	input.Rol = strings.ToLower(input.Rol)
	if !rolesValidos[input.Rol] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Rol no permitido"})
		return
	}

	var existente models.Usuario
	if err := db.Where("nombre = ?", input.Nombre).First(&existente); err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "El usuario ya existe"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Clave), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al encriptar la clave"})
		return
	}

	usuario := models.Usuario{
		Nombre: input.Nombre,
		Clave:  string(hash),
		Rol:    input.Rol,
	}

	if err := db.Create(&usuario); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear el usuario"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"mensaje": "Usuario creado correctamente"})
}
