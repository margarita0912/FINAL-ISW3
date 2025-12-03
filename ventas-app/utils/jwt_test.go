package utils

import (
	"os"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	// Setup
	os.Setenv("JWT_SECRET", "test-secret-key")
	secret = []byte(os.Getenv("JWT_SECRET"))

	t.Run("genera token válido", func(t *testing.T) {
		token, err := GenerateToken(1, "admin")

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("token contiene claims correctos", func(t *testing.T) {
		userID := uint(123)
		rol := "vendedor"

		tokenStr, err := GenerateToken(userID, rol)
		assert.NoError(t, err)

		// Parsear el token para verificar claims
		token, err := ParseToken(tokenStr)
		assert.NoError(t, err)
		assert.True(t, token.Valid)

		claims, ok := token.Claims.(jwt.MapClaims)
		assert.True(t, ok)
		assert.Equal(t, float64(userID), claims["user_id"])
		assert.Equal(t, rol, claims["rol"])
		assert.NotNil(t, claims["exp"])
	})
}

func TestParseToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	secret = []byte(os.Getenv("JWT_SECRET"))

	t.Run("parsea token válido correctamente", func(t *testing.T) {
		// Generar token primero
		tokenStr, _ := GenerateToken(1, "admin")

		token, err := ParseToken(tokenStr)

		assert.NoError(t, err)
		assert.NotNil(t, token)
		assert.True(t, token.Valid)
	})

	t.Run("falla con token inválido", func(t *testing.T) {
		token, err := ParseToken("token-invalido")

		assert.Error(t, err)
		assert.Nil(t, token)
	})

	t.Run("falla con token malformado", func(t *testing.T) {
		token, err := ParseToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.malformed")

		assert.Error(t, err)
		assert.Nil(t, token)
	})
}
