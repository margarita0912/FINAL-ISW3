package controllers

import (
	"net/http"
	"ventas-app/database"
	"ventas-app/models"

	"github.com/gin-gonic/gin"
)

func RegistrarCompra(c *gin.Context) {
	db := database.GetDB(c) // ← selecciona la base según el entorno

	var compra models.Compra
	if err := c.ShouldBindJSON(&compra); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// Validación adicional: producto_id y cantidad deben ser mayores que 0
	if compra.ProductoID <= 0 || compra.Cantidad <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	var producto models.Producto
	if err := db.First(&producto, compra.ProductoID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Producto no encontrado"})
		return
	}

	producto.Stock += compra.Cantidad

	if err := db.Save(&producto); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el stock"})
		return
	}

	if err := db.Create(&compra); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al registrar la compra"})
		return
	}

	c.JSON(http.StatusCreated, compra)
}
