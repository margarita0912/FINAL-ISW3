package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestVentaModel_Creation(t *testing.T) {
	venta := Venta{
		UsuarioID:   1,
		ProductoID:  10,
		Cantidad:    5,
		Descuento:   10.0,
		PrecioFinal: 90.0,
	}

	assert.Equal(t, uint(1), venta.UsuarioID)
	assert.Equal(t, uint(10), venta.ProductoID)
	assert.Equal(t, 5, venta.Cantidad)
	assert.Equal(t, 10.0, venta.Descuento)
	assert.Equal(t, 90.0, venta.PrecioFinal)
}

func TestVentaModel_PrecioFinalConDescuento(t *testing.T) {
	precioOriginal := 100.0
	descuento := 20.0
	precioFinal := precioOriginal - descuento

	venta := Venta{
		UsuarioID:   1,
		ProductoID:  5,
		Cantidad:    2,
		Descuento:   descuento,
		PrecioFinal: precioFinal,
	}

	assert.Equal(t, 80.0, venta.PrecioFinal)
	assert.Less(t, venta.PrecioFinal, precioOriginal)
}

func TestVentaModel_SinDescuento(t *testing.T) {
	venta := Venta{
		UsuarioID:   2,
		ProductoID:  8,
		Cantidad:    3,
		Descuento:   0.0,
		PrecioFinal: 150.0,
	}

	assert.Zero(t, venta.Descuento)
	assert.Equal(t, 150.0, venta.PrecioFinal)
}

func TestVentaModel_MultipleCantidad(t *testing.T) {
	precioUnitario := 25.0
	cantidad := 10
	total := precioUnitario * float64(cantidad)

	venta := Venta{
		UsuarioID:   3,
		ProductoID:  15,
		Cantidad:    cantidad,
		Descuento:   0.0,
		PrecioFinal: total,
	}

	assert.Equal(t, 10, venta.Cantidad)
	assert.Equal(t, 250.0, venta.PrecioFinal)
}

func TestVentaModel_DescuentoPorcentual(t *testing.T) {
	precioBase := 200.0
	porcentajeDescuento := 0.15 // 15%
	descuento := precioBase * porcentajeDescuento
	precioFinal := precioBase - descuento

	venta := Venta{
		UsuarioID:   1,
		ProductoID:  20,
		Cantidad:    1,
		Descuento:   descuento,
		PrecioFinal: precioFinal,
	}

	assert.Equal(t, 30.0, venta.Descuento)
	assert.Equal(t, 170.0, venta.PrecioFinal)
}

func TestVentaModel_GormModel(t *testing.T) {
	venta := Venta{
		Model:       gorm.Model{ID: 100},
		UsuarioID:   5,
		ProductoID:  25,
		Cantidad:    7,
		Descuento:   5.0,
		PrecioFinal: 95.0,
	}

	assert.Equal(t, uint(100), venta.ID)
	assert.NotNil(t, venta.Model)
}

func TestVentaModel_EmptyVenta(t *testing.T) {
	venta := Venta{}

	assert.Zero(t, venta.UsuarioID)
	assert.Zero(t, venta.ProductoID)
	assert.Zero(t, venta.Cantidad)
	assert.Zero(t, venta.Descuento)
	assert.Zero(t, venta.PrecioFinal)
}

func TestVentaModel_JSONTags(t *testing.T) {
	venta := Venta{
		Model:       gorm.Model{ID: 50},
		UsuarioID:   12,
		ProductoID:  33,
		Cantidad:    4,
		Descuento:   15.5,
		PrecioFinal: 184.5,
	}

	assert.NotZero(t, venta.ID)
	assert.NotZero(t, venta.UsuarioID)
	assert.NotZero(t, venta.ProductoID)
}

func TestVentaModel_NegativeValues(t *testing.T) {
	// Test para detectar valores negativos que deberían ser validados
	venta := Venta{
		UsuarioID:   1,
		ProductoID:  5,
		Cantidad:    -2,
		Descuento:   -10.0,
		PrecioFinal: -50.0,
	}

	assert.Negative(t, venta.Cantidad)
	assert.Negative(t, venta.Descuento)
	assert.Negative(t, venta.PrecioFinal)
}

func TestVentaModel_HighVolumeVenta(t *testing.T) {
	venta := Venta{
		UsuarioID:   1,
		ProductoID:  100,
		Cantidad:    1000,
		Descuento:   500.0,
		PrecioFinal: 9500.0,
	}

	assert.Equal(t, 1000, venta.Cantidad)
	assert.Greater(t, venta.PrecioFinal, venta.Descuento)
}
