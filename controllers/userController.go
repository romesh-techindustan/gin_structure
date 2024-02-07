package controllers

import (
	"crypto/rand"
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

func CreateUser(c *gin.Context) {
	user, exist := c.Get("user")
	data := user.(models.User)
	if !exist {
		c.JSON(http.StatusForbidden, "Unable to changed password")
	} else {
		if data.IsSuperuser {
			var body struct {
				Name     string
				Email    string
				Password string
			}
			if c.Bind(&body) != nil {
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
			}
			password := generatePassword(6)
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
			body.Password = string(hashedPassword)
			user = models.User{ID: uuid.New().String(), Name: body.Name, Email: body.Email, Password: body.Password}
			result = database.Db.Create(&user)
			fmt.Println("password", password)
			if result.Error != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Unable to create user!",
				})
			}
			c.JSON(http.StatusOK, gin.H{
				"password": password,
				"result":   user,
			})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "You are not authorized to create user",
			})
		}

	}

}

func generatePassword(length int) string {
	b := make([]byte, length+2)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[2 : length+2]
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func UserLogin(c *gin.Context) {
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

	otp := GenerateTOTPCode(os.Getenv("OTP_SECRET_KEY"))
	database.Db.Model(&user).Update("otp", otp)
	emailParams := &services.EmailParams{
		To:   user.Email,
		Code: otp,
	}
	services.Send2FAOTP(*emailParams)
	c.JSON(http.StatusOK, gin.H{
		"response":     "Login Successfull",
		"access_token": accessToken,
		"otp":          otp,
	})
}

func ResetPassword(c *gin.Context) {
	var email string
	c.Bind(&email)
	url := "http://localhost:3000/changedpwd"
	emailParams := &services.EmailParams{
		To:               email,
		PasswordResetURL: url,
	}

	services.SendResetPasswordEmail(*emailParams)
}

type ChangePasswordRequest struct {
	ID              string `json:"id"`
	Password        string `json:"password1"`
	ConfirmPassword string `json:"password"`
}

func ChangePassword(c *gin.Context) {
	var pass ChangePasswordRequest
	c.BindJSON(&pass)
	if pass.Password != pass.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Different password",
		})
	} else {

		// Update the user's password
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(pass.Password), 10)

		// Save the updated user to the database
		database.Db.Model(&models.User{}).Where("id = ?", pass.ID).Update("password", string(hashedPassword))
		c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
	}

}

func GetAllUsers(c *gin.Context) {
	user, exist := c.Get("user")
	data := user.(models.User)
	if !exist {
		c.JSON(http.StatusForbidden, "Unable to changed password")
	} else {
		if data.IsSuperuser {
			var Users []models.User
			result := database.Db.Find(&Users)
			if result.Error != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Users not found",
				})
			} else {
				c.JSON(200, gin.H{
					"users": Users,
				})
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "This function is enabled to superadmin only",
			})
		}

	}
}

func DeleteUser(c *gin.Context) {
	user, exist := c.Get("user")
	data := user.(models.User)
	if !exist {
		c.JSON(http.StatusForbidden, "Unable to changed password")
	} else {
		if data.IsSuperuser {
			id := c.Param("id")
			var user models.User
			result := database.Db.Where("id = ?", id).Delete(&user)
			if result.Error != nil {
				c.JSON(http.StatusForbidden, gin.H{
					"error": "User not found",
				})
			} else {
				c.JSON(200, gin.H{
					"message": "User Deleted Successfully",
				})
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "This function is enabled to superadmin only",
			})
		}
	}
}

func GetUserDetail(c *gin.Context) {
	user, exist := c.Get("user")
	_ = user.(models.User)
	if !exist {
		c.JSON(http.StatusForbidden, "Unable to changed password")
	} else {
		id := c.Param("id")
		var user models.User
		result := database.Db.Where("id = ?", id).First(&user)
		if result.Error != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "User not found",
			})
		} else {
			c.JSON(200, gin.H{
				"user": user,
			})
		}
	}
}

func Logout(c *gin.Context) {
	// Clear the cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", "", -1, "", "", false, true)

	// Redirect to the home page
	c.JSON(200, gin.H{
		"message": "User Log Out",
	})

}
