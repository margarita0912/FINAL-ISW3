package controllers

import (
	"net/http"
	"ventas-app/database"
	"ventas-app/models"

	"github.com/gin-gonic/gin"
)

func CrearProducto(c *gin.Context) {
	db := database.GetDB(c) // ← selecciona la base según el entorno

	var p models.Producto
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	if err := db.Create(&p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar"})
		return
	}

	c.JSON(http.StatusCreated, p)
}

func ListarProductos(c *gin.Context) {
	db := database.GetDB(c) // ← selecciona la base según el entorno

	var productos []models.Producto
	if err := db.Find(&productos); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al listar productos"})
		return
	}

	c.JSON(http.StatusOK, productos)
}
