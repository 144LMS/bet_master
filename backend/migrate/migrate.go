package main

import (
	"github.com/144LMS/bet_master/initializers"
	"github.com/144LMS/bet_master/models"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectToDB()
}

func main() {
	initializers.DB.AutoMigrate(&models.User{}, &models.Wallet{}, &models.Transaction{}, &models.Admin{}, &models.Bet{}, &models.Match{})
}
