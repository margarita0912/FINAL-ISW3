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

func TestRegistrarCompra_InvalidBody(t *testing.T) {
	database.GetDB = func(c *gin.Context) database.DBHandler {
		return &mocks.MockDB{ShouldErr: false}
	}

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/compras", RegistrarCompra)

	body := `{"producto_id": 0, "cantidad": 0}`
	req, _ := http.NewRequest("POST", "/compras", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

// Test: JSON malformado
func TestRegistrarCompra_MalformedJSON(t *testing.T) {
	database.GetDB = func(c *gin.Context) database.DBHandler {
		return &mocks.MockDB{ShouldErr: false}
	}

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/compras", RegistrarCompra)

	body := `{"producto_id": 1, "cantidad":}`  // JSON incompleto
	req, _ := http.NewRequest("POST", "/compras", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	
	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "Datos inválidos")
}

// Test: Producto no encontrado
func TestRegistrarCompra_ProductNotFound(t *testing.T) {
	mock := &mocks.MockDB{
		Productos: []models.Producto{}, // Sin productos
		ShouldErr: false,
	}

	database.GetDB = func(c *gin.Context) database.DBHandler {
		return mock
	}

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/compras", RegistrarCompra)

	body := `{"producto_id": 999, "cantidad": 5}`  // Producto inexistente
	req, _ := http.NewRequest("POST", "/compras", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
	
	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "Producto no encontrado")
}

// Test: Error al guardar producto
func TestRegistrarCompra_SaveError(t *testing.T) {
	prod := models.Producto{Model: gorm.Model{ID: 1}, Nombre: "P1", Precio: 10.0, Stock: 5}
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
	router.POST("/compras", RegistrarCompra)

	body := `{"producto_id": 1, "cantidad": 3}`
	req, _ := http.NewRequest("POST", "/compras", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	
	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "Error al actualizar el stock")
}

// Test: Error al crear compra
func TestRegistrarCompra_CreateCompraError(t *testing.T) {
	prod := models.Producto{Model: gorm.Model{ID: 1}, Nombre: "P1", Precio: 10.0, Stock: 5}
	mock := &mocks.MockDB{
		Productos:  []models.Producto{prod},
		ShouldErr:  false,
		FailCreate: true, // Simular error al crear compra
	}

	database.GetDB = func(c *gin.Context) database.DBHandler {
		return mock
	}

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/compras", RegistrarCompra)

	body := `{"producto_id": 1, "cantidad": 3}`
	req, _ := http.NewRequest("POST", "/compras", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	
	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "Error al registrar la compra")
}
