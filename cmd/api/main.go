package main

import (
	"github.com/144LMS/bet_master.git/initializers"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnv()
}

func main() {
	r := gin.Default()

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "statusOK",
		})
	})

	r.Run()
}
