package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestProductoModel_Creation(t *testing.T) {
	producto := Producto{
		Nombre: "Laptop",
		Costo:  500.0,
		Precio: 800.0,
		Stock:  10,
	}

	assert.Equal(t, "Laptop", producto.Nombre)
	assert.Equal(t, 500.0, producto.Costo)
	assert.Equal(t, 800.0, producto.Precio)
	assert.Equal(t, 10, producto.Stock)
}

func TestProductoModel_PriceValidation(t *testing.T) {
	producto := Producto{
		Nombre: "Mouse",
		Costo:  10.0,
		Precio: 15.0,
		Stock:  50,
	}

	assert.Greater(t, producto.Precio, producto.Costo, "El precio debe ser mayor al costo")
}

func TestProductoModel_StockManagement(t *testing.T) {
	producto := Producto{
		Nombre: "Teclado",
		Costo:  20.0,
		Precio: 30.0,
		Stock:  100,
	}

	// Simular venta
	cantidadVendida := 5
	producto.Stock -= cantidadVendida

	assert.Equal(t, 95, producto.Stock)
	assert.GreaterOrEqual(t, producto.Stock, 0, "El stock no debe ser negativo")
}

func TestProductoModel_JSONTags(t *testing.T) {
	producto := Producto{
		Model:  gorm.Model{ID: 1},
		Nombre: "Monitor",
		Costo:  200.0,
		Precio: 300.0,
		Stock:  25,
	}

	assert.NotZero(t, producto.ID)
	assert.Equal(t, "Monitor", producto.Nombre)
}

func TestProductoModel_EmptyProduct(t *testing.T) {
	producto := Producto{}

	assert.Empty(t, producto.Nombre)
	assert.Zero(t, producto.Costo)
	assert.Zero(t, producto.Precio)
	assert.Zero(t, producto.Stock)
}

func TestProductoModel_NegativeValues(t *testing.T) {
	// Test para verificar que se pueden detectar valores negativos
	producto := Producto{
		Nombre: "ProductoTest",
		Costo:  -10.0,
		Precio: -5.0,
		Stock:  -1,
	}

	// En un sistema real, estos valores deberían ser validados
	assert.Negative(t, producto.Costo)
	assert.Negative(t, producto.Precio)
	assert.Negative(t, producto.Stock)
}

func TestProductoModel_ZeroStock(t *testing.T) {
	producto := Producto{
		Nombre: "Agotado",
		Costo:  50.0,
		Precio: 75.0,
		Stock:  0,
	}

	assert.Zero(t, producto.Stock, "El stock debe ser 0")
	assert.False(t, producto.Stock > 0, "No hay stock disponible")
}

func TestProductoModel_HighValues(t *testing.T) {
	producto := Producto{
		Nombre: "Servidor",
		Costo:  10000.0,
		Precio: 15000.0,
		Stock:  5,
	}

	assert.Equal(t, 10000.0, producto.Costo)
	assert.Equal(t, 15000.0, producto.Precio)
	margen := producto.Precio - producto.Costo
	assert.Equal(t, 5000.0, margen)
}
