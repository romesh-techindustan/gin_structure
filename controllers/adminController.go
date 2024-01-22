package controllers

import (
	"fmt"
	"net/http"
	"os"
	"zucora/backend/config"
	"zucora/backend/database"
	"zucora/backend/models"
	"zucora/backend/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CreateSuperAdmin(c *gin.Context) {
	var body struct {
		Name     string
		Email    string
		Password string
	}
	if c.BindJSON(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read Body",
		})
	}
	var user models.User
	result := database.Db.Where(&models.User{Email: body.Email}).First(&user)
	if result.RowsAffected != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"eror": "User already exists",
		})
	} else {
		password := generatePassword(6)
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
		body.Password = string(hashedPassword)
		user = models.User{ID: uuid.New().String(), Name: body.Name, Email: body.Email, Password: body.Password, IsSuperuser: true}
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
