package main

import (
	"github.com/144LMS/bet_master.git/initializers"
	"github.com/144LMS/bet_master.git/internal/user"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectToDB()
}

func main() {
	r := gin.Default()

	r.GET("/user/:id", user.GetUserController)

	r.POST("/user", user.CreateUserController)

	r.DELETE("/user/:id", user.DeleteUserController)

	r.PUT("/user/:id", user.UpdateUserController)

	r.Run()
}
