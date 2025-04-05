package middleware

import (
	"net/http"
	"strings"

	"github.com/144LMS/bet_master/internal/auth"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authService *auth.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		userID, userRole, err := authService.ValidateTokens(strings.TrimPrefix(tokenString, "Bearer "))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		ctx.Set("userID", userID)
		ctx.Set("role", userRole)
		ctx.Next()
	}
}
