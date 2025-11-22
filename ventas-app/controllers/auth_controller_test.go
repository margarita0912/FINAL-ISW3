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
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func TestLogin_InvalidBody(t *testing.T) {
	database.GetDB = func(c *gin.Context) database.DBHandler {
		return &mocks.MockDB{ShouldErr: false}
	}

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/login", Login)

	body := `{"nombre": "", "clave": ""}`
	req, _ := http.NewRequest("POST", "/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

// Test: JSON malformado
func TestLogin_MalformedJSON(t *testing.T) {
	database.GetDB = func(c *gin.Context) database.DBHandler {
		return &mocks.MockDB{ShouldErr: false}
	}

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/login", Login)

	body := `{"nombre": "test", "clave":`  // JSON incompleto
	req, _ := http.NewRequest("POST", "/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

// Test: Login exitoso
func TestLogin_Success(t *testing.T) {
	// Crear clave hasheada para el test
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	
	mock := &mocks.MockDB{
		Usuarios: []models.Usuario{
			{
				Model:  gorm.Model{ID: 1},
				Nombre: "testuser",
				Clave:  string(hashedPassword),
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
	router.POST("/login", Login)

	body := `{"nombre": "testuser", "clave": "password123"}`
	req, _ := http.NewRequest("POST", "/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "token")
	assert.Equal(t, "vendedor", response["rol"])
}

// Test: Usuario no encontrado
func TestLogin_UserNotFound(t *testing.T) {
	mock := &mocks.MockDB{
		Usuarios:  []models.Usuario{}, // Sin usuarios
		ShouldErr: false,
	}

	database.GetDB = func(c *gin.Context) database.DBHandler {
		return mock
	}

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/login", Login)

	body := `{"nombre": "inexistente", "clave": "password123"}`
	req, _ := http.NewRequest("POST", "/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

// Test: Clave incorrecta
func TestLogin_WrongPassword(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
	
	mock := &mocks.MockDB{
		Usuarios: []models.Usuario{
			{
				Model:  gorm.Model{ID: 1},
				Nombre: "testuser",
				Clave:  string(hashedPassword),
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
	router.POST("/login", Login)

	body := `{"nombre": "testuser", "clave": "wrongpassword"}`
	req, _ := http.NewRequest("POST", "/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}
