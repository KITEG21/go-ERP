package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"user_api/internal/auth"
	"user_api/internal/common"
)

func AuthMiddleware(jwtService *auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, common.APIError{Message: "Missing authorization header", Code: "401"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, common.APIError{Message: "Invalid authorization header", Code: "401"})
			c.Abort()
			return
		}

		claims, err := jwtService.ValidateToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, common.APIError{Message: "Invalid token", Code: "401"})
			c.Abort()
			return
		}

		c.Set("user_id", claims["id"])
		c.Set("user_email", claims["email"])
		c.Next()
	}
}
