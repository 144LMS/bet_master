package user

import (
	"net/http"

	"github.com/144LMS/bet_master/internal/auth"
	"github.com/144LMS/bet_master/models"
	"github.com/gin-gonic/gin"
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
	userID := ctx.Param("id")

	user, err := c.userService.getUserService(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	ctx.JSON(http.StatusOK, user)
}

func (c *Controller) GetUserWithWalletController(ctx *gin.Context) {
	userID := ctx.Param("id")

	user, wallet, err := c.userService.getUserWithWalletService(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	responce := gin.H{
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

	ctx.JSON(http.StatusOK, responce)
}

func (c *Controller) Registration(ctx *gin.Context) {
	var user models.User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdUser, createdWallet, err := c.userService.CreateUserWithWalletService(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	responce := gin.H{
		"user": gin.H{
			"id":       createdUser.ID,
			"username": createdUser.Username,
			"email":    createdUser.Email,
		},
		"wallet": gin.H{
			"id":      createdWallet.ID,
			"balance": createdWallet.Balance,
		},
	}

	ctx.JSON(http.StatusOK, responce)

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
	var userData struct {
		Email    string `json:"email" gorm:"unique"`
		Password string `json:"password" gorm:"not null"`
	}

	if err := ctx.ShouldBindJSON(&userData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.userService.Autenticate(userData.Email, userData.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	accessToken, refreshToken, err := c.authService.GenerateTokens(user.ID, user.Role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
