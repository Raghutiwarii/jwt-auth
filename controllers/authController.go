package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"jwt-auth/database"

	"jwt-auth/models"
)

func OnBoardingUser(c *gin.Context) {
	var (
		user = models.User{}
	)

	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_code":  "A001",
			"description": "Bad format",
		})
	}

	var existingUser models.User
	err = database.DB.Where("email = ?", user.Email).First(&existingUser).Error
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email already exists",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
	}

	newUser := models.User{
		Name:     user.Name,
		Email:    user.Email,
		Address:  user.Address,
		Password: string(hashedPassword),
	}

	err = database.DB.Create(&newUser).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user",
		})
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User is created",
	})
}
