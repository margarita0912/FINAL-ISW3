package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSetup_HealthEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	Setup(router)

	req, _ := http.NewRequest("GET", "/healthz", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "OK", resp.Body.String())
}

func TestSetup_HealthEndpointHEAD(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	Setup(router)

	req, _ := http.NewRequest("HEAD", "/healthz", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestSetup_LoginRouteExists(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	Setup(router)

	// Verificar que la ruta /login existe (aunque falle sin datos)
	req, _ := http.NewRequest("POST", "/login", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Debería retornar algo (no 404)
	assert.NotEqual(t, http.StatusNotFound, resp.Code)
}

func TestSetup_UsuariosRouteExists(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	Setup(router)

	req, _ := http.NewRequest("POST", "/usuarios", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.NotEqual(t, http.StatusNotFound, resp.Code)
}

func TestSetup_ProductosGetRouteExists(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	Setup(router)

	req, _ := http.NewRequest("GET", "/productos", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.NotEqual(t, http.StatusNotFound, resp.Code)
}

func TestSetup_ProductosHeadRouteExists(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	Setup(router)

	req, _ := http.NewRequest("HEAD", "/productos", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.NotEqual(t, http.StatusNotFound, resp.Code)
}

func TestSetup_ProductosPostRouteExists(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	Setup(router)

	req, _ := http.NewRequest("POST", "/productos", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.NotEqual(t, http.StatusNotFound, resp.Code)
}

func TestSetup_ComprasRouteExists(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	Setup(router)

	req, _ := http.NewRequest("POST", "/compras", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.NotEqual(t, http.StatusNotFound, resp.Code)
}

func TestSetup_VentasRouteExists(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	Setup(router)

	req, _ := http.NewRequest("POST", "/ventas", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.NotEqual(t, http.StatusNotFound, resp.Code)
}

func TestSetup_InvalidRouteReturns404(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	Setup(router)

	req, _ := http.NewRequest("GET", "/invalid/route", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestSetup_AllRoutesRegistered(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	Setup(router)

	routes := router.Routes()

	// Debería haber al menos 8 rutas registradas
	assert.GreaterOrEqual(t, len(routes), 8)
}
