package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"kopibang/domain/dto"
	"kopibang/internal/tokenutil"
)

func AuthMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse(http.StatusUnauthorized, "Missing or invalid token", nil))
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := tokenutil.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse(http.StatusUnauthorized, "Token expired or invalid", err.Error()))
			c.Abort()
			return
		}

		// Jika rute butuh role spesifik (misal "barista"), tolak jika tidak sesuai
		if requiredRole != "" && claims.Role != requiredRole {
			c.JSON(http.StatusForbidden, dto.ErrorResponse(http.StatusForbidden, "You don't have permission to access this resource", nil))
			c.Abort()
			return
		}

		// Set data penting ke dalam context Gin agar bisa dipakai oleh Controller
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}