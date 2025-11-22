package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"ventas-app/database"
	"ventas-app/mocks"
	"ventas-app/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestRegistrarVenta_InvalidBody(t *testing.T) {
	database.GetDB = func(c *gin.Context) database.DBHandler {
		return &mocks.MockDB{ShouldErr: false}
	}

	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.POST("/ventas", RegistrarVenta)

	// Arrange: sin datos de producto ni cantidad
	body := `{"producto_id": 0, "cantidad": 0}`
	req, _ := http.NewRequest("POST", "/ventas", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	// Act
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, resp.Code, "Debe devolver 400 si los datos son inválidos")
}

// Test: JSON malformado
func TestRegistrarVenta_MalformedJSON(t *testing.T) {
	database.GetDB = func(c *gin.Context) database.DBHandler {
		return &mocks.MockDB{ShouldErr: false}
	}

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/ventas", RegistrarVenta)

	body := `{"producto_id": 1, "cantidad":}`  // JSON incompleto
	req, _ := http.NewRequest("POST", "/ventas", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	
	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "Datos inválidos")
}

// Test: Producto no encontrado
func TestRegistrarVenta_ProductNotFound(t *testing.T) {
	mock := &mocks.MockDB{
		Productos: []models.Producto{}, // Sin productos
		ShouldErr: false,
	}

	database.GetDB = func(c *gin.Context) database.DBHandler {
		return mock
	}

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/ventas", RegistrarVenta)

	body := `{"producto_id": 999, "cantidad": 2}`  // Producto inexistente
	req, _ := http.NewRequest("POST", "/ventas", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
	
	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "Producto no encontrado")
}

// Test: Stock insuficiente
func TestRegistrarVenta_InsufficientStock(t *testing.T) {
	prod := models.Producto{Model: gorm.Model{ID: 1}, Nombre: "P1", Precio: 20.0, Stock: 5} // Solo 5 en stock
	mock := &mocks.MockDB{
		Productos: []models.Producto{prod},
		ShouldErr: false,
	}

	database.GetDB = func(c *gin.Context) database.DBHandler {
		return mock
	}

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/ventas", RegistrarVenta)

	body := `{"producto_id": 1, "cantidad": 10}`  // Intentar vender más de lo disponible
	req, _ := http.NewRequest("POST", "/ventas", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	
	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "Stock insuficiente")
}

// Test: Error al guardar producto
func TestRegistrarVenta_SaveProductError(t *testing.T) {
	prod := models.Producto{Model: gorm.Model{ID: 1}, Nombre: "P1", Precio: 20.0, Stock: 10}
	mock := &mocks.MockDB{
		Productos: []models.Producto{prod},
		ShouldErr: false,
		FailSave:  true, // Simular error al guardar
	}

	database.GetDB = func(c *gin.Context) database.DBHandler {
		return mock
	}

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/ventas", RegistrarVenta)

	body := `{"producto_id": 1, "cantidad": 2}`
	req, _ := http.NewRequest("POST", "/ventas", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	
	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "Error al actualizar el stock")
}

// Test: Error al crear venta
func TestRegistrarVenta_CreateVentaError(t *testing.T) {
	prod := models.Producto{Model: gorm.Model{ID: 1}, Nombre: "P1", Precio: 20.0, Stock: 10}
	mock := &mocks.MockDB{
		Productos:  []models.Producto{prod},
		ShouldErr:  false,
		FailCreate: true, // Simular error al crear venta
	}

	database.GetDB = func(c *gin.Context) database.DBHandler {
		return mock
	}

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/ventas", RegistrarVenta)

	body := `{"producto_id": 1, "cantidad": 2}`
	req, _ := http.NewRequest("POST", "/ventas", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	
	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "Error al registrar la venta")
}
