package main

import (
	"github.com/144LMS/bet_master.git/initializers"
	"github.com/144LMS/bet_master.git/models"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectToDB()
}

func main() {
	initializers.DB.AutoMigrate(&models.User{}, &models.Admin{}, &models.Bet{})
}
