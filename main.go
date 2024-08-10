package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/register", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	})
	r.Run()
}
