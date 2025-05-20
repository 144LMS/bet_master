package main

import (
	"os"
	"time"

	//"github.com/144LMS/bet_master/initializers"
	"github.com/144LMS/bet_master/initializers"
	"github.com/144LMS/bet_master/internal/admin"
	"github.com/144LMS/bet_master/internal/auth"
	"github.com/144LMS/bet_master/internal/user"
	"github.com/144LMS/bet_master/internal/wallet"
	"github.com/144LMS/bet_master/middleware"
	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectToDB()
}

func main() {

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // Cookie, В Postman "Send cookies with requests" в настройках
		MaxAge:           12 * time.Hour,
	}))

	r.Static("/", "D:/vsCode/bet_master/web/dist")

	userRepo := user.NewUserRepository(initializers.DB)
	walletRepo := wallet.NewWalletRepository(initializers.DB)
	matchRepo := admin.NewMatchRepository(initializers.DB)
	betRepo := admin.NewBetRepository(initializers.DB)
	adminRepo := admin.NewAdminRepository(initializers.DB)

	userService := user.NewUserService(userRepo)
	authService := auth.NewAuthService(os.Getenv("JWT_SECRET"))
	walletService := wallet.NewWalletService(walletRepo)
	matchService := admin.NewMatchService(*matchRepo)
	betService := admin.NewBetService(walletRepo, matchRepo, betRepo)
	adminService := admin.NewAdminService(*adminRepo)

	userController := user.NewUserController(userService, authService)
	walletController := wallet.NewWalletController(walletService)
	adminController := admin.NewAdminController(
		betService,
		matchService,
		adminService,
		os.Getenv("JWT_SECRET"),
	)

	api := r.Group("/")
	{
		api.POST("/register", userController.Registration)
		api.POST("/login", userController.Login)
		api.POST("/admin/login", adminController.AdminLogin)
		api.POST("/admin/register", adminController.CreateAdmin)
	}

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware(authService))
	{
		protected.GET("/user/:id", userController.GetUserWithWalletController)
		protected.PUT("/user/:id", userController.UpdateUserController)
		protected.DELETE("/user/:id", userController.DeleteUserController)

		//protected.GET("/wallets/:user_id", walletController.GetWallet)
		protected.POST("/wallets/:user_id/deposit", walletController.Deposit)
		protected.POST("/wallets/:user_id/withdraw", walletController.Withdraw)
		protected.GET("/wallets/:user_id/balance", walletController.GetBalance)
		protected.GET("/wallets/:user_id/transactions", walletController.GetTransactions)
	}

	adminRoutes := r.Group("/admin")
	adminRoutes.Use(middleware.AuthMiddleware(authService))
	adminRoutes.OPTIONS("/matches", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.AbortWithStatus(204)
	})
	adminRoutes.Use(middleware.AdminMiddleware())
	{
		adminRoutes.GET("/dashboard", adminController.AdminDashboard)
		adminRoutes.GET("/users", userController.GetUserController)
		adminRoutes.GET("/matches", adminController.GetAllMatches)
		adminRoutes.POST("/matches", adminController.CreateMatch)
		adminRoutes.DELETE("/matches/:id", adminController.DeleteMatch)
		adminRoutes.POST("/matches/:id/settle", adminController.SettleMatch)

	}

	r.Run()
}
