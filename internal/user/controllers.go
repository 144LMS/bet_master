package user

import (
	"net/http"

	"github.com/144LMS/bet_master.git/initializers"
	"github.com/144LMS/bet_master.git/models"
	"github.com/gin-gonic/gin"
)

func GetUserController(ctx *gin.Context) {
	userID := ctx.Param("id")

	var user models.User

	if err := initializers.DB.First(&user, userID).Error; err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	userResponce := models.UserResponse{
		ID:        int(user.ID),
		Username:  user.Username,
		Email:     user.Email,
		Balance:   user.Balance,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt.Time,
		Role:      user.Role,
	}

	ctx.JSON(http.StatusOK, userResponce)
}

func CreateUserController(ctx *gin.Context) {
	var user models.User

	if err := ctx.Bind(&user); err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	result := initializers.DB.Create(&user)

	if err := result.Error; err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(200, user)

}

func DeleteUserController(ctx *gin.Context) {
	var deleteUserReq models.DeleteUserRequest
	deleteUserReq.ID = ctx.Param("id")

	if err := initializers.DB.Where("id = ?", deleteUserReq.ID).Delete(&models.User{}).Error; err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
	}
	ctx.JSON(200, deleteUserReq.ID)
}

func UpdateUserController(ctx *gin.Context) {
	var userUpdate models.UpdateUserRequest
	var user models.User
	userID := ctx.Param("id")

	if err := ctx.ShouldBindJSON(&userUpdate); err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := initializers.DB.First(&user, userID).Error; err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	user.Username = userUpdate.Username
	user.Email = userUpdate.Email
	user.Password = userUpdate.Password

	if err := initializers.DB.Save(&user).Error; err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"user was update": userUpdate,
	})
}
