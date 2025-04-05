package main

import (
	"os"

	"github.com/144LMS/bet_master/initializers"
	"github.com/144LMS/bet_master/internal/auth"
	"github.com/144LMS/bet_master/internal/user"
	"github.com/144LMS/bet_master/middleware"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectToDB()
}

func main() {
	r := gin.Default()

	userRepo := user.NewUserRepository(initializers.DB)

	userService := user.NewUserService(userRepo)
	authService := auth.NewAuthService(os.Getenv("JWT_SECRET"))

	userController := user.NewUserController(userService, authService)

	public := r.Group("/")
	{
		public.POST("/login", userController.Login)
		public.POST("/register", userController.Registration)
	}

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware(authService))
	{
		protected.GET("/users/:id", userController.GetUserController)
		protected.PUT("/users/:id", userController.UpdateUserController)
		protected.DELETE("/users/:id", userController.DeleteUserController)
	}
	r.Run()
}

/*
userRoutes := router.Group("/user")
{
    userRoutes.POST("/login", userController.Login)
    userRoutes.POST("/register", userController.Register)
    // ...
}

adminRoutes := router.Group("/admin")
{
    adminRoutes.POST("/login", adminController.AdminLogin) // Отдельный endpoint!
    adminRoutes.GET("/users", adminController.GetAllUsers)
    // ...
}
*/
