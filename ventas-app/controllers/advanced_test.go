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

// Test avanzado: listar productos (happy path) y verificar cuerpo JSON
func TestListarProductos_HappyPath(t *testing.T) {
	// Arrange
	mock := &mocks.MockDB{
		Productos: []models.Producto{{Model: gorm.Model{ID: 1}, Nombre: "ProductoX", Precio: 9.99, Stock: 10}},
		ShouldErr: false,
	}
	database.GetDB = func(c *gin.Context) database.DBHandler {
		return mock
	}

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/productos", ListarProductos)

	req, _ := http.NewRequest("GET", "/productos", nil)
	resp := httptest.NewRecorder()

	// Act
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusOK, resp.Code)
	var prods []models.Producto
	err := json.Unmarshal(resp.Body.Bytes(), &prods)
	assert.NoError(t, err)
	assert.Len(t, prods, 1)
	assert.Equal(t, "ProductoX", prods[0].Nombre)
}

// Test avanzado: listar productos cuando la DB falla
func TestListarProductos_DBError(t *testing.T) {
	// Arrange
	mock := &mocks.MockDB{ShouldErr: true}
	database.GetDB = func(c *gin.Context) database.DBHandler { return mock }

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/productos", ListarProductos)

	req, _ := http.NewRequest("GET", "/productos", nil)
	resp := httptest.NewRecorder()

	// Act
	router.ServeHTTP(resp, req)

	// Assert: en caso de error en Find esperamos 500
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
}

// Test de excepción: login cuando la DB devuelve error (no hacer panic)
func TestLogin_DBError_NoPanic(t *testing.T) {
	// Arrange
	mock := &mocks.MockDB{ShouldErr: true}
	database.GetDB = func(c *gin.Context) database.DBHandler { return mock }

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/login", Login)

	body := `{"nombre": "inexistente", "clave": "x"}`
	req, _ := http.NewRequest("POST", "/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	// Act
	router.ServeHTTP(resp, req)

	// Assert: controller debe manejar el error y devolver 401 (usuario no encontrado)
	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

// Test: Crear producto cuando la DB falla al crear -> 500
func TestCrearProducto_CreateError(t *testing.T) {
	// Arrange
	mock := &mocks.MockDB{FailCreate: true}
	database.GetDB = func(c *gin.Context) database.DBHandler { return mock }

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/productos", CrearProducto)

	body := `{"nombre":"X","costo":5.0,"precio":9.0,"stock":10}`
	req, _ := http.NewRequest("POST", "/productos", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	// Act
	router.ServeHTTP(resp, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
}

// Test happy path: Registrar compra actualiza stock y crea compra
func TestRegistrarCompra_HappyPath(t *testing.T) {
	// Arrange: producto con stock 5
	prod := models.Producto{Model: gorm.Model{ID: 1}, Nombre: "P1", Precio: 10.0, Stock: 5}
	mock := &mocks.MockDB{Productos: []models.Producto{prod}}
	database.GetDB = func(c *gin.Context) database.DBHandler { return mock }

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/compras", RegistrarCompra)

	body := `{"producto_id":1,"cantidad":3}`
	req, _ := http.NewRequest("POST", "/compras", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	// Act
	router.ServeHTTP(resp, req)

	// Assert: creado y stock incrementado a 8
	assert.Equal(t, http.StatusCreated, resp.Code)
	if len(mock.Productos) > 0 {
		assert.Equal(t, 8, mock.Productos[0].Stock)
	}
}

// Test happy path: Registrar venta calcula precio y decrementa stock
func TestRegistrarVenta_HappyPath(t *testing.T) {
	// Arrange: producto con stock 10, precio 20
	prod := models.Producto{Model: gorm.Model{ID: 2}, Nombre: "P2", Precio: 20.0, Stock: 10}
	mock := &mocks.MockDB{Productos: []models.Producto{prod}}
	database.GetDB = func(c *gin.Context) database.DBHandler { return mock }

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/ventas", RegistrarVenta)

	body := `{"producto_id":2,"cantidad":2,"descuento":10}`
	req, _ := http.NewRequest("POST", "/ventas", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	// Act
	router.ServeHTTP(resp, req)

	// Assert: creado y stock decrementado a 8
	assert.Equal(t, http.StatusCreated, resp.Code)
	if len(mock.Productos) > 0 {
		assert.Equal(t, 8, mock.Productos[0].Stock)
	}
}
