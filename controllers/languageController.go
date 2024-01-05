package controllers
import (
	"github.com/gin-gonic/gin"
	"zucora/backend/models"
	"net/http"
)


func ChangeLanguage(c *gin.Context){
	val,exist:=c.Get("user")
	val = val.(models.User)
	if !exist{
		c.JSON(http.StatusBadRequest, gin.H{
			"error":"Invalid User",
		})
	}
	var language string
	err := c.Bind(&language)
	if err!= nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error":"Invalid Language",
		})
	}
}