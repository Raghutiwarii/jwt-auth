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
	r.POST("/login", controllers.Login)
	r.GET("/users", controllers.GetAllUsers)
	r.POST("/register", controllers.OnBoardingUser)
	r.GET("/user/:user_id", controllers.GetUser)
	r.Run()
}
