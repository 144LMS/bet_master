package main

import (
	"os"

	"github.com/144LMS/bet_master/initializers"
	"github.com/144LMS/bet_master/internal/auth"
	"github.com/144LMS/bet_master/internal/user"
	"github.com/144LMS/bet_master/internal/wallet"
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
	walletRepo := wallet.NewWalletRepository(initializers.DB)

	userService := user.NewUserService(userRepo)
	authService := auth.NewAuthService(os.Getenv("JWT_SECRET"))
	walletService := wallet.NewWalletService(walletRepo)

	userController := user.NewUserController(userService, authService)
	walletController := wallet.NewWalletController(walletService)

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

		protected.GET("/wallets/:user_id", walletController.GetWallet)
		protected.POST("/wallets/:user_id/deposit", walletController.Deposit)
		protected.POST("/wallets/:user_id/withdraw", walletController.Withdraw)
		protected.GET("/wallets/:user_id/balance", walletController.GetBalance)
		protected.GET("/wallets/:user_id/transactions", walletController.GetTransactions)
	}

	adminRoutes := r.Group("/admin")
	adminRoutes.Use(middleware.AuthMiddleware(authService))
	adminRoutes.Use(middleware.AdminMiddleware())
	{
		/*
		   adminRoutes.GET("/users", adminController.GetAllUsers)
		   adminRoutes.POST("/users", adminController.CreateUser)
		   adminRoutes.DELETE("/users/:id", adminController.ForceDeleteUser)

		   adminRoutes.POST("/wallets/:user_id/force-deposit", adminController.ForceDeposit)
		   adminRoutes.POST("/wallets/:user_id/force-withdraw", adminController.ForceWithdraw)
		   adminRoutes.POST("/wallets/:user_id/set-balance", adminController.SetBalance)
		   adminRoutes.GET("/wallets", adminController.GetAllWallets)
		   adminRoutes.GET("/transactions", adminController.GetAllTransactions)
		*/
	}

	r.Run()
}

/*
func main() {
    r := gin.Default()

    // Инициализация репозиториев
    userRepo := user.NewUserRepository(initializers.DB)
    walletRepo := repositories.NewWalletRepository(initializers.DB)

    // Инициализация сервисов
    userService := user.NewUserService(userRepo)
    authService := auth.NewAuthService(os.Getenv("JWT_SECRET"))
    walletService := services.NewWalletService(walletRepo)

    // Инициализация контроллеров
    userController := user.NewUserController(userService, authService)
    walletController := controllers.NewWalletController(walletService)
    adminController := admin.NewAdminController(userService, walletService)

    // Публичные роуты (не требуют авторизации)
    public := r.Group("/")
    {
        public.POST("/login", userController.Login)
        public.POST("/register", userController.Registration)
    }

    // Защищенные роуты (требуют авторизации)
    protected := r.Group("/")
    protected.Use(middleware.AuthMiddleware(authService))
    {
        // Пользовательские эндпоинты
        protected.GET("/users/:id", userController.GetUserController)
        protected.PUT("/users/:id", userController.UpdateUserController)
        protected.DELETE("/users/:id", userController.DeleteUserController)

        // Эндпоинты для работы с балансом
        protected.GET("/wallets/:user_id", walletController.GetWallet)
        protected.POST("/wallets/:user_id/deposit", walletController.Deposit)
        protected.POST("/wallets/:user_id/withdraw", walletController.Withdraw)
        protected.GET("/wallets/:user_id/balance", walletController.GetBalance)
        protected.GET("/wallets/:user_id/transactions", walletController.GetTransactions)
    }

    // Админские роуты (требуют авторизации и прав администратора)
    adminRoutes := r.Group("/admin")
    adminRoutes.Use(middleware.AuthMiddleware(authService))
    adminRoutes.Use(middleware.AdminMiddleware())
    {
        // Управление пользователями
        adminRoutes.GET("/users", adminController.GetAllUsers)
        adminRoutes.POST("/users", adminController.CreateUser)
        adminRoutes.DELETE("/users/:id", adminController.ForceDeleteUser)

        // Управление балансами
        adminRoutes.POST("/wallets/:user_id/force-deposit", adminController.ForceDeposit)
        adminRoutes.POST("/wallets/:user_id/force-withdraw", adminController.ForceWithdraw)
        adminRoutes.POST("/wallets/:user_id/set-balance", adminController.SetBalance)
        adminRoutes.GET("/wallets", adminController.GetAllWallets)
        adminRoutes.GET("/transactions", adminController.GetAllTransactions)
    }

    r.Run()
}
*/
