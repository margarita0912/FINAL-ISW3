package controllers

import (
	"net/http"
	"ventas-app/database"
	"ventas-app/models"

	"github.com/gin-gonic/gin"
)

func RegistrarVenta(c *gin.Context) {
	db := database.GetDB(c) // ← selecciona la base según el entorno

	var venta models.Venta
	if err := c.ShouldBindJSON(&venta); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// Validación adicional: producto_id y cantidad deben ser mayores que 0
	if venta.ProductoID <= 0 || venta.Cantidad <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	var producto models.Producto
	if err := db.First(&producto, venta.ProductoID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Producto no encontrado"})
		return
	}

	if producto.Stock < venta.Cantidad {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Stock insuficiente"})
		return
	}

	precioUnitario := producto.Precio
	total := float64(venta.Cantidad) * precioUnitario
	totalConDescuento := total * (1 - venta.Descuento/100)
	totalConIVA := totalConDescuento * 1.21

	venta.PrecioFinal = totalConIVA
	producto.Stock -= venta.Cantidad

	if err := db.Save(&producto); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el stock"})
		return
	}

	if err := db.Create(&venta); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al registrar la venta"})
		return
	}

	c.JSON(http.StatusCreated, venta)
}
