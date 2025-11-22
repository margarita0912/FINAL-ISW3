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

func TestCrearUsuario_InvalidData(t *testing.T) {

	database.GetDB = func(c *gin.Context) database.DBHandler {
		return &mocks.MockDB{ShouldErr: false}
	}

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/usuarios", CrearUsuario)

	body := `{"nombre": "", "clave": "", "rol": ""}`
	req, _ := http.NewRequest("POST", "/usuarios", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

// Test: JSON malformado
func TestCrearUsuario_MalformedJSON(t *testing.T) {
	database.GetDB = func(c *gin.Context) database.DBHandler {
		return &mocks.MockDB{ShouldErr: false}
	}

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/usuarios", CrearUsuario)

	body := `{"nombre": "test", "clave":}`  // JSON incompleto
	req, _ := http.NewRequest("POST", "/usuarios", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

// Test: Rol inválido
func TestCrearUsuario_InvalidRole(t *testing.T) {
	database.GetDB = func(c *gin.Context) database.DBHandler {
		return &mocks.MockDB{ShouldErr: false}
	}

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/usuarios", CrearUsuario)

	body := `{"nombre": "testuser", "clave": "password123", "rol": "admin"}`
	req, _ := http.NewRequest("POST", "/usuarios", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	
	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "Rol no permitido")
}

// Test: Usuario ya existente
func TestCrearUsuario_UserExists(t *testing.T) {
	mock := &mocks.MockDB{
		Usuarios: []models.Usuario{
			{
				Model:  gorm.Model{ID: 1},
				Nombre: "existing",
				Clave:  "hashedpassword",
				Rol:    "vendedor",
			},
		},
		ShouldErr: false,
	}

	database.GetDB = func(c *gin.Context) database.DBHandler {
		return mock
	}

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/usuarios", CrearUsuario)

	body := `{"nombre": "existing", "clave": "password123", "rol": "vendedor"}`
	req, _ := http.NewRequest("POST", "/usuarios", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusConflict, resp.Code)
	
	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "El usuario ya existe")
}

// Test: Error al crear usuario en la BD
func TestCrearUsuario_DatabaseError(t *testing.T) {
	mock := &mocks.MockDB{
		Usuarios:   []models.Usuario{}, // Sin usuarios existentes
		ShouldErr:  false,
		FailCreate: true, // Simular error al crear
	}

	database.GetDB = func(c *gin.Context) database.DBHandler {
		return mock
	}

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/usuarios", CrearUsuario)

	body := `{"nombre": "newuser", "clave": "password123", "rol": "vendedor"}`
	req, _ := http.NewRequest("POST", "/usuarios", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	
	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "Error al crear el usuario")
}

// Test: Crear usuario exitoso
func TestCrearUsuario_Success(t *testing.T) {
	mock := &mocks.MockDB{
		Usuarios:  []models.Usuario{}, // Sin usuarios existentes
		ShouldErr: false,
	}

	database.GetDB = func(c *gin.Context) database.DBHandler {
		return mock
	}

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/usuarios", CrearUsuario)

	body := `{"nombre": "newuser", "clave": "password123", "rol": "VENDEDOR"}`  // Rol en mayúscula para probar conversión
	req, _ := http.NewRequest("POST", "/usuarios", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
	
	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Contains(t, response["mensaje"], "Usuario creado correctamente")
}
