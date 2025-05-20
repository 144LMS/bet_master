package middleware

import (
	"net/http"
	"strings"

	"github.com/144LMS/bet_master/internal/auth"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authService *auth.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Пробуем получить токен из кук сначала
		tokenString, err := ctx.Cookie("access_token")
		if err != nil {
			// Если нет в куках, пробуем из заголовка Authorization
			authHeader := ctx.GetHeader("Authorization")
			if authHeader == "" {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Требуется авторизация: отсутствует токен в куках или заголовке Authorization",
				})
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Неверный формат заголовка Authorization",
				})
				return
			}

			tokenString = parts[1]
			if tokenString == "" {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Токен не предоставлен",
				})
				return
			}
		}

		userID, userRole, err := authService.ValidateTokens(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Неверный токен",
				"details": err.Error(),
			})
			return
		}

		ctx.Set("userID", userID)
		ctx.Set("role", userRole)
		ctx.Next()
	}
}
