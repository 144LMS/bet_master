package user

import (
	"fmt"
	"net/http"

	"github.com/144LMS/bet_master/internal/auth"
	"github.com/144LMS/bet_master/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Controller struct {
	userService *UserService
	authService *auth.AuthService
}

func NewUserController(userService *UserService, authService *auth.AuthService) *Controller {
	return &Controller{
		userService: userService,
		authService: authService}
}

func (c *Controller) GetUserController(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	user, err := c.userService.repo.GetUserRepository(fmt.Sprintf("%v", userID))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

func (c *Controller) GetUserWithWalletController(ctx *gin.Context) {
	userID := ctx.Param("id")

	user, wallet, err := c.userService.getUserWithWalletService(userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	response := gin.H{
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
		"wallet": gin.H{
			"id":      wallet.ID,
			"balance": wallet.Balance,
		},
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *Controller) Registration(ctx *gin.Context) {
	var request struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные", "details": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании пользователя"})
		return
	}

	user := models.User{
		Username: request.Username,
		Email:    request.Email,
		Password: string(hashedPassword),
	}

	createdUser, createdWallet, err := c.userService.CreateUserWithWalletService(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"input": gin.H{
				"username": request.Username,
				"email":    request.Email,
			},
		})
		return
	}

	accessToken, refreshToken, err := c.authService.GenerateTokens(createdUser.ID, createdUser.Role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации токенов"})
		return
	}

	ctx.SetCookie("access_token",
		accessToken,
		3600,
		"/",
		"",
		false,
		true,
	)

	ctx.SetCookie(
		"refresh_token",
		refreshToken,
		86400,
		"/",
		"",
		false,
		true,
	)

	ctx.JSON(http.StatusOK, gin.H{
		//"access_token":  accessToken,
		//"refresh_token": refreshToken,
		"user": gin.H{
			"id":       createdUser.ID,
			"username": createdUser.Username,
			"email":    createdUser.Email,
		},
		"wallet": gin.H{
			"id":      createdWallet.ID,
			"balance": createdWallet.Balance,
		},
	})
}

func (c *Controller) DeleteUserController(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.userService.DeleteUserService(id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, id)
}

func (c *Controller) UpdateUserController(ctx *gin.Context) {
	var req models.UpdateUserRequest
	id := ctx.Param("id")

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.userService.UpdateUserService(id, req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "user updated"})
}

func (c *Controller) Login(ctx *gin.Context) {
	var credentials struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.userService.Authenticate(credentials.Email, credentials.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	accessToken, refreshToken, err := c.authService.GenerateTokens(user.ID, user.Role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate tokens"})
		return
	}

	ctx.SetCookie("access_token",
		accessToken,
		3600,
		"/",
		"",
		false,
		true,
	)

	ctx.SetCookie(
		"refresh_token",
		refreshToken,
		86400,
		"/",
		"",
		false,
		true,
	)

	ctx.JSON(http.StatusOK, gin.H{
		//"access_token":  accessToken,
		//"refresh_token": refreshToken,
		"user_id": user.ID,
	})
}
