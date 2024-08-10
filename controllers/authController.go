package controllers

import (
	"jwt-auth/database"
	"jwt-auth/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// OnBoardingUser handles user registration
func OnBoardingUser(c *gin.Context) {
	var user models.User

	// Bind JSON request body to user struct
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_code":  "A001",
			"description": "Bad format",
		})
		return
	}

	// Check if the email already exists
	var existingUser models.User
	if err := database.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email already exists",
		})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	// Create new user
	newUser := models.User{
		Name:     user.Name,
		Email:    user.Email,
		Address:  user.Address,
		Password: string(hashedPassword),
	}

	// Save new user to the database
	if err := database.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User is created",
	})
}

// GetAllUsers retrieves all users from the database
func GetAllUsers(c *gin.Context) {
	var users []models.User // Use a slice to store multiple users

	// Fetch all users from the database
	if err := database.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve users",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}
