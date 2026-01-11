package handler

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/user-service/internal/core/domain"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("jwt_token")
		if err != nil {
			// fallback: check Authorization header
			authHeader := c.GetHeader("Authorization")
			tokenString, _ = strings.CutPrefix(authHeader, "Bearer ")
		}

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			c.Abort()
			return
		}

		// parse token
		secret := os.Getenv("JWT_SECRET")
		token, err := jwt.ParseWithClaims(tokenString, &domain.JWTClaims{}, func(token *jwt.Token) (any, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			c.Abort()
			return
		}

		// extract claims
		claims, ok := token.Claims.(*domain.JWTClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token claims",
			})
			c.Abort()
			return
		}

		// store user claims in context for downstream handlers
		c.Set("user", claims)
		c.Next()
	}
}

// ensures that the user has at least one of the required roles
func RoleMiddleware(requiredRoles ...domain.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "no user in context",
			})
			c.Abort()
			return
		}

		claims, ok := user.(*domain.JWTClaims)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "invalid user claims",
			})
			c.Abort()
			return
		}

		authorized := false
		for _, role := range requiredRoles {
			if claims.Role == role {
				authorized = true
				break
			}
		}

		if !authorized {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "forbidden",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}