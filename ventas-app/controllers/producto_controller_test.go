package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"ventas-app/database"
	"ventas-app/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestListarProductos_OK(t *testing.T) {
	database.GetDB = func(c *gin.Context) database.DBHandler {
		return &mocks.MockDB{ShouldErr: false}
	}

	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.GET("/productos", ListarProductos)

	req, _ := http.NewRequest("GET", "/productos", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

// Test: JSON malformado al crear producto
func TestCrearProducto_MalformedJSON(t *testing.T) {
	database.GetDB = func(c *gin.Context) database.DBHandler {
		return &mocks.MockDB{ShouldErr: false}
	}

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/productos", CrearProducto)

	body := `{"nombre": "Producto", "precio":}`  // JSON incompleto
	req, _ := http.NewRequest("POST", "/productos", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	
	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "Datos inválidos")
}

// Test: Crear producto exitoso
func TestCrearProducto_Success(t *testing.T) {
	mock := &mocks.MockDB{ShouldErr: false}

	database.GetDB = func(c *gin.Context) database.DBHandler {
		return mock
	}

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/productos", CrearProducto)

	body := `{"nombre": "Producto Test", "precio": 15.99, "stock": 100}`
	req, _ := http.NewRequest("POST", "/productos", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
}
