package controllers

import (
	"fmt"
	"jwt-auth/database"
	"jwt-auth/models"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var jwtSecret = []byte(os.Getenv("SECRET"))

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

func GetUser(c *gin.Context) {
	// Extract the user ID from URL parameters
	userIDStr := c.Param("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User ID is required",
		})
		return
	}

	// Convert the user ID to an integer
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	// Fetch the user from the database
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to retrieve user",
			})
		}
		return
	}

	// Return the user details
	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// Login handles user authentication
func Login(c *gin.Context) {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Bind JSON request body to loginData struct
	if err := c.BindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_code":  "A002",
			"description": "Invalid request format",
		})
		return
	}

	// Fetch user by email
	var user models.User
	if err := database.DB.Where("email = ?", loginData.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email",
		})
		return
	}

	userInput, err := bcrypt.GenerateFromPassword([]byte(loginData.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error hashing password:", err)
		return
	}

	// Compare the provided password with the stored hashed password
	if err := bcrypt.CompareHashAndPassword(userInput, []byte(loginData.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid password",
		})
		return
	}

	// Generate a JWT token
	token, err := generateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	// Return the JWT token
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}

// generateToken generates a JWT token for the given user ID
func generateToken(user models.User) (string, error) {
	// Define the expiration time for the token
	expirationTime := time.Now().Add(5 * time.Minute)

	// Create the claims
	claims := jwt.MapClaims{
		"Subject":   strconv.Itoa(int(user.ID)),
		"ExpiresAt": jwt.NewNumericDate(expirationTime),
		"Email":     user.Email,
	}

	// Create the token with the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with a secret key
	secretKey := []byte(jwtSecret)
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func GetUserProfile(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Fetch user data from the database
	var user models.User // Assuming User is your model for user data
	if err := database.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Return user profile data
	c.JSON(http.StatusOK, gin.H{
		"id":      user.ID,
		"name":    user.Name,
		"email":   user.Email,
		"address": user.Address,
	})
}
