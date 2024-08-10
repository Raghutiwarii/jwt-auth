package main

import (
	"fmt"
	"net/http"

	"jwt-auth/database"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	dbSuccess, err := database.ConnectDB()
	fmt.Println(dbSuccess)
	if(err!=nil){
		panic("could not connect to db")
	}
	r.GET("/register", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	})
	r.Run()
}
