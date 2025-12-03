package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCompraModel_Creation(t *testing.T) {
	compra := Compra{
		UsuarioID:  1,
		ProductoID: 10,
		Cantidad:   50,
		CostoUnit:  25.5,
	}

	assert.Equal(t, uint(1), compra.UsuarioID)
	assert.Equal(t, uint(10), compra.ProductoID)
	assert.Equal(t, 50, compra.Cantidad)
	assert.Equal(t, 25.5, compra.CostoUnit)
}

func TestCompraModel_CostoTotal(t *testing.T) {
	compra := Compra{
		UsuarioID:  2,
		ProductoID: 15,
		Cantidad:   10,
		CostoUnit:  100.0,
	}

	costoTotal := float64(compra.Cantidad) * compra.CostoUnit
	assert.Equal(t, 1000.0, costoTotal)
}

func TestCompraModel_SingleUnit(t *testing.T) {
	compra := Compra{
		UsuarioID:  3,
		ProductoID: 20,
		Cantidad:   1,
		CostoUnit:  45.75,
	}

	assert.Equal(t, 1, compra.Cantidad)
	costoTotal := float64(compra.Cantidad) * compra.CostoUnit
	assert.Equal(t, compra.CostoUnit, costoTotal)
}

func TestCompraModel_BulkPurchase(t *testing.T) {
	compra := Compra{
		UsuarioID:  1,
		ProductoID: 5,
		Cantidad:   500,
		CostoUnit:  10.0,
	}

	costoTotal := float64(compra.Cantidad) * compra.CostoUnit
	assert.Equal(t, 5000.0, costoTotal)
	assert.GreaterOrEqual(t, compra.Cantidad, 100, "Compra al por mayor")
}

func TestCompraModel_LowCostUnit(t *testing.T) {
	compra := Compra{
		UsuarioID:  4,
		ProductoID: 8,
		Cantidad:   100,
		CostoUnit:  0.50,
	}

	assert.Equal(t, 0.50, compra.CostoUnit)
	costoTotal := float64(compra.Cantidad) * compra.CostoUnit
	assert.Equal(t, 50.0, costoTotal)
}

func TestCompraModel_HighCostUnit(t *testing.T) {
	compra := Compra{
		UsuarioID:  5,
		ProductoID: 100,
		Cantidad:   5,
		CostoUnit:  2500.0,
	}

	assert.Equal(t, 2500.0, compra.CostoUnit)
	costoTotal := float64(compra.Cantidad) * compra.CostoUnit
	assert.Equal(t, 12500.0, costoTotal)
}

func TestCompraModel_GormModel(t *testing.T) {
	compra := Compra{
		Model:      gorm.Model{ID: 75},
		UsuarioID:  10,
		ProductoID: 30,
		Cantidad:   25,
		CostoUnit:  15.0,
	}

	assert.Equal(t, uint(75), compra.ID)
	assert.NotNil(t, compra.Model)
}

func TestCompraModel_EmptyCompra(t *testing.T) {
	compra := Compra{}

	assert.Zero(t, compra.UsuarioID)
	assert.Zero(t, compra.ProductoID)
	assert.Zero(t, compra.Cantidad)
	assert.Zero(t, compra.CostoUnit)
}

func TestCompraModel_JSONTags(t *testing.T) {
	compra := Compra{
		Model:      gorm.Model{ID: 50},
		UsuarioID:  15,
		ProductoID: 40,
		Cantidad:   75,
		CostoUnit:  12.99,
	}

	assert.NotZero(t, compra.ID)
	assert.NotZero(t, compra.UsuarioID)
	assert.NotZero(t, compra.ProductoID)
	assert.NotZero(t, compra.Cantidad)
	assert.NotZero(t, compra.CostoUnit)
}

func TestCompraModel_NegativeValues(t *testing.T) {
	// Test para detectar valores negativos que deberían ser validados
	compra := Compra{
		UsuarioID:  1,
		ProductoID: 5,
		Cantidad:   -10,
		CostoUnit:  -5.0,
	}

	assert.Negative(t, compra.Cantidad)
	assert.Negative(t, compra.CostoUnit)
}

func TestCompraModel_ZeroCost(t *testing.T) {
	compra := Compra{
		UsuarioID:  2,
		ProductoID: 10,
		Cantidad:   50,
		CostoUnit:  0.0,
	}

	assert.Zero(t, compra.CostoUnit)
	costoTotal := float64(compra.Cantidad) * compra.CostoUnit
	assert.Zero(t, costoTotal, "Costo total debe ser 0")
}

func TestCompraModel_DecimalCost(t *testing.T) {
	compra := Compra{
		UsuarioID:  3,
		ProductoID: 12,
		Cantidad:   7,
		CostoUnit:  9.99,
	}

	assert.Equal(t, 9.99, compra.CostoUnit)
	costoTotal := float64(compra.Cantidad) * compra.CostoUnit
	assert.InDelta(t, 69.93, costoTotal, 0.01)
}

func TestCompraModel_MultipleProducts(t *testing.T) {
	compras := []Compra{
		{UsuarioID: 1, ProductoID: 1, Cantidad: 10, CostoUnit: 5.0},
		{UsuarioID: 1, ProductoID: 2, Cantidad: 20, CostoUnit: 3.0},
		{UsuarioID: 1, ProductoID: 3, Cantidad: 15, CostoUnit: 7.0},
	}

	var totalCompra float64
	for _, c := range compras {
		totalCompra += float64(c.Cantidad) * c.CostoUnit
	}

	assert.Equal(t, 3, len(compras))
	assert.Equal(t, 215.0, totalCompra)
}
