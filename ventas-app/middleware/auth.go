package middleware

import (
	"net/http"
	"strings"
	"ventas-app/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthRequired(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token faltante"})
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := utils.ParseToken(tokenStr)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		rol := claims["rol"].(string)

		// Validar rol
		for _, r := range roles {
			if r == rol {
				c.Set("user_id", claims["user_id"])
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Acceso denegado"})
		c.Abort()
	}
}
