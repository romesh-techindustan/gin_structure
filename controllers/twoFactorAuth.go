package controllers

import (
	"net/http"
	"time"
	"zucora/backend/models"

	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
)

func VerifyTwoFA(c *gin.Context) {
	var otp struct {
		OTP string
	}
	c.Bind(&otp)
	user, exist := c.Get("user")
	data := user.(models.User)
	if !exist {
		c.JSON(http.StatusForbidden, "Unable to changed password")
	} else {
		if otp.OTP != data.OTP {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Invalid OTP .",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"message": "2FA validation successful",
			})
		}
	}
}

func GenerateTOTPCode(secretKey string) string {
	val, _ := totp.GenerateCode(secretKey, time.Now())
	return val
}
