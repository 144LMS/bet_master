package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthController struct {
	service *AuthService
}

func NewAuthController(service *AuthService) *AuthController {
	return &AuthController{service: service}
}

// POST /auth/refresh
func (c *AuthController) RefreshToken(ctx *gin.Context) {
	var request struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	token, err := jwt.Parse(request.RefreshToken, func(t *jwt.Token) (interface{}, error) {
		return []byte("refresh_secret"), nil
	})

	if err != nil || !token.Valid {
		ctx.JSON(401, gin.H{"error": "Invalid refresh token"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		ctx.JSON(401, gin.H{"error": "Invalid token claims"})
		return
	}

	userID, ok := claims["sub"].(uint)
	if !ok {
		ctx.JSON(401, gin.H{"error": "Invalid user ID in token"})
		return
	}

	userRole, ok := claims["role"].(string)
	if !ok {
		ctx.JSON(401, gin.H{"error": "Invalid user ID in token"})
		return
	}

	newAccess, newRefresh, _ := c.service.GenerateTokens(userID, userRole)

	ctx.JSON(200, gin.H{
		"access_token":  newAccess,
		"refresh_token": newRefresh,
	})
}
