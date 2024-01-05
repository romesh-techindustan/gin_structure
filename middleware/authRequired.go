package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"
	"zucora/backend/database"
	"zucora/backend/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthRequired(c *gin.Context) {

	tokenString := c.GetHeader("Authorization")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["sub"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusFound, "token not valid")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			fmt.Println("Token Expired")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token expired",
			})
		}
		str := claims["userID"].(string)
		var user models.User
		result := database.Db.Where(&models.User{ID: str}).Find(&user)
		if result.RowsAffected == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "User not found",
			})
		}
		fmt.Println(user)
		c.Set("user", user)
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Header("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		c.Next()
	} else {
		c.JSON(http.StatusUnauthorized, err)
	}
}
