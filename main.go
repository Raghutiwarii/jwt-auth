package main

import (
	"jwt-auth/controllers"
	"jwt-auth/database"
	"jwt-auth/initializer"
	"jwt-auth/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	initializer.LoadEnvVariables()
	r := gin.Default()

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Database connection
	_, err := database.ConnectDB()
	if err != nil {
		panic("could not connect to db")
	}

	// Public route for login
	r.POST("/login", controllers.Login)
	r.POST("/register", controllers.OnBoardingUser)

	// Secure routes with JWT authentication middleware
	secured := r.Group("/")
	secured.Use(middleware.AuthMiddleware())

	secured.GET("/users", controllers.GetAllUsers)
	secured.GET("/user/:user_id", controllers.GetUser)
	secured.GET("/user/profile", controllers.GetUserProfile)

	// Start the server
	r.Run(":8080")
}
