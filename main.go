package main

import (
	"jwt-auth/initializer"
	"jwt-auth/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	initializer.LoadEnvVariables()
	r := gin.Default()

	_, err := database.ConnectDB()
	if err != nil {
		panic("could not connect to db")
	}
	r.GET("/register", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	})
	r.Run()
}
