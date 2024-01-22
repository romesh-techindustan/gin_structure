package middleware

import (
	"zucora/backend/services"

	"github.com/gin-gonic/gin"
)

func AuthRequired(c *gin.Context) {
	tokenParams := &services.TokenParams{
		AuthToken: c.GetHeader("Authorization"),
	}
	user, status, err := services.VerifyAccessToken(*tokenParams)

	if err != nil {
		c.JSON(status, err)
	}

	c.Set("user", user)
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Header("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
	c.Next()

}
