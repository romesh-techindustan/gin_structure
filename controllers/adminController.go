package controllers

import (
	"fmt"
	"net/http"
	"os"
	"zucora/backend/config"
	"zucora/backend/database"
	"zucora/backend/models"
	"zucora/backend/requests"
	"zucora/backend/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CreateSuperAdmin(c *gin.Context) {
	var input requests.CreateAdmin
	var user models.User

	if c.BindJSON(&input) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
	}
	result := database.Db.Where(&models.User{Email: input.Email}).First(&user)
	if result.RowsAffected != 0 {
		c.JSON(http.StatusConflict, gin.H{
			"error": "User already exists",
		})
	}
	password := generatePassword(6)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	input.Password = string(hashedPassword)
	user = models.User{
		ID:          uuid.New().String(),
		Name:        input.Name,
		Email:       input.Email,
		Password:    input.Password,
		IsSuperuser: true,
	}
	result = database.Db.Create(&user)
	fmt.Println("password", password)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unable to create user!",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"result": user,
	})
}

func SuperAdminLogin(c *gin.Context) {
	var loginRequest LoginRequest

	// Bind the JSON body to the LoginRequest struct
	if err := c.Bind(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user models.User
	// Access the data from the JSON body
	email := loginRequest.Email
	password := loginRequest.Password
	fmt.Println(email)
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please enter an email address"})
		return
	}
	result := database.Db.First(&user, "email = ?", email)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect Email or Password !"})
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "incorrect password",
		})
		return
	}
	// Generate access and refresh tokens
	tokenParams := &services.TokenParams{
		TokenExpiration: config.GetConfig().TokenExpiration,
		JWTSecret:       []byte(os.Getenv("JWT_SECRET_KEY")),
		UserID:          user.ID,
	}
	accessToken, err := services.GenerateAccessToken(*tokenParams)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"response":     "Login Successfull",
		"access_token": accessToken,
	})
}
