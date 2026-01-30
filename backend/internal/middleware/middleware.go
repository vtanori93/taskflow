package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/taskflow/backend/internal/models"
	"github.com/taskflow/backend/internal/utils/jwt"
)

// AuthMiddleware middleware para validar JWT
func AuthMiddleware(jwtManager *jwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener el token del header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Data:       nil,
				StatusCode: http.StatusUnauthorized,
				Message:    "Token no proporcionado",
				Error:      true,
			})
			c.Abort()
			return
		}

		// Extraer el token del header Bearer
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Data:       nil,
				StatusCode: http.StatusUnauthorized,
				Message:    "Formato de token inválido",
				Error:      true,
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validar el token
		claims, err := jwtManager.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Data:       nil,
				StatusCode: http.StatusUnauthorized,
				Message:    "Token inválido o expirado",
				Error:      true,
			})
			c.Abort()
			return
		}

		// Guardar datos en contexto
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)

		c.Next()
	}
}

// CORSMiddleware middleware para CORS - permite peticiones desde cualquier origen
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}

		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "false")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// ErrorHandlingMiddleware middleware para manejar panics
func ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				errMsg := fmt.Sprintf("%v", err)
				fmt.Printf("🔴 PANIC RECOVERED: %v\n", err)
				c.JSON(http.StatusInternalServerError, models.APIResponse{
					Data:       nil,
					StatusCode: http.StatusInternalServerError,
					Message:    "Error interno del servidor: " + errMsg,
					Error:      true,
				})
			}
		}()
		c.Next()
	}
}
