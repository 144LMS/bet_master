package main

import (
	"os"
	"time"

	"github.com/144LMS/bet_master/initializers"
	"github.com/144LMS/bet_master/internal/admin"
	"github.com/144LMS/bet_master/internal/auth"
	"github.com/144LMS/bet_master/internal/bet"
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
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	userRepo := user.NewUserRepository(initializers.DB)
	walletRepo := wallet.NewWalletRepository(initializers.DB)
	matchRepo := admin.NewMatchRepository(initializers.DB)
	betRepoForAdmin := admin.NewBetRepository(initializers.DB)
	adminRepo := admin.NewAdminRepository(initializers.DB)
	betRepo := bet.NewBetRepository(initializers.DB)

	userService := user.NewUserService(userRepo)
	authService := auth.NewAuthService(os.Getenv("JWT_SECRET"))
	walletService := wallet.NewWalletService(walletRepo)
	matchService := admin.NewMatchService(*matchRepo)
	betServiceForAdmin := admin.NewBetService(walletRepo, matchRepo, betRepoForAdmin)
	adminService := admin.NewAdminService(*adminRepo)
	betService := bet.NewBetService(*betRepo)

	userController := user.NewUserController(userService, authService)
	walletController := wallet.NewWalletController(walletService)
	adminController := admin.NewAdminController(
		betServiceForAdmin,
		matchService,
		adminService,
		os.Getenv("JWT_SECRET"),
	)
	betConroller := bet.NewBetController(*betService, *walletService)

	api := r.Group("/api")
	{

		api.POST("/register", userController.Registration)

		api.POST("/login", userController.Login)
		api.POST("/admin/login", adminController.AdminLogin)
		api.POST("/admin/register", adminController.CreateAdmin)
	}

	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware(authService))
	{

		protected.GET("/user/me", userController.GetUserController)

		protected.GET("/user/:id", userController.GetUserWithWalletController)

		protected.PUT("/user/:id", userController.UpdateUserController)

		protected.DELETE("/user/:id", userController.DeleteUserController)

		protected.GET("/wallets", walletController.GetWallet)

		protected.POST("/wallets/deposit", walletController.Deposit)

		protected.POST("/wallets/withdraw", walletController.Withdraw)
		protected.GET("/wallets/balance", walletController.GetBalance)
		protected.GET("/wallets/transactions", walletController.GetTransactions)
		protected.GET("/bet/matches", adminController.GetAllMatches)

		protected.POST("/bet/placeBet", betConroller.PlaceBet)

	}

	adminRoutes := r.Group("/api/admin")
	adminRoutes.Use(middleware.AuthMiddleware(authService))
	adminRoutes.OPTIONS("/matches", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.AbortWithStatus(204)
	})
	adminRoutes.Use(middleware.AdminMiddleware())
	{

		adminRoutes.GET("/dashboard", adminController.AdminDashboard)

		adminRoutes.GET("/matches", adminController.GetAllMatches)

		adminRoutes.POST("/matches", adminController.CreateMatch)

		adminRoutes.DELETE("/matches/:id", adminController.DeleteMatch)

		adminRoutes.POST("/matches/:id/settle", adminController.SettleMatch)

	}

	r.Run()
}
