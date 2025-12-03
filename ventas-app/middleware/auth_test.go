package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"ventas-app/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthRequired(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	os.Setenv("JWT_SECRET", "test-secret-key")

	t.Run("rechaza request sin token", func(t *testing.T) {
		router := gin.New()
		router.GET("/protected", AuthRequired("admin"), func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "ok"})
		})

		req := httptest.NewRequest("GET", "/protected", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
		assert.Contains(t, resp.Body.String(), "Token faltante")
	})

	t.Run("rechaza token sin Bearer prefix", func(t *testing.T) {
		router := gin.New()
		router.GET("/protected", AuthRequired("admin"), func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "ok"})
		})

		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "invalid-token")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})

	t.Run("rechaza token inválido", func(t *testing.T) {
		router := gin.New()
		router.GET("/protected", AuthRequired("admin"), func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "ok"})
		})

		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer token-invalido")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
		assert.Contains(t, resp.Body.String(), "Token inválido")
	})

	t.Run("acepta token válido con rol correcto", func(t *testing.T) {
		token, _ := utils.GenerateToken(1, "admin")

		router := gin.New()
		router.GET("/protected", AuthRequired("admin"), func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "ok"})
		})

		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("rechaza token válido con rol incorrecto", func(t *testing.T) {
		token, _ := utils.GenerateToken(1, "vendedor")

		router := gin.New()
		router.GET("/protected", AuthRequired("admin"), func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "ok"})
		})

		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusForbidden, resp.Code)
		assert.Contains(t, resp.Body.String(), "Acceso denegado")
	})

	t.Run("acepta múltiples roles válidos", func(t *testing.T) {
		token, _ := utils.GenerateToken(1, "vendedor")

		router := gin.New()
		router.GET("/protected", AuthRequired("admin", "vendedor"), func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "ok"})
		})

		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
	})
}
