package main

import (
	"jwt-auth/controllers"
	"jwt-auth/database"
	"jwt-auth/initializer"

	"github.com/gin-gonic/gin"
)

func main() {
	initializer.LoadEnvVariables()
	r := gin.Default()

	_, err := database.ConnectDB()
	if err != nil {
		panic("could not connect to db")
	}
	r.POST("/register", controllers.OnBoardingUser)
	r.Run()
}
