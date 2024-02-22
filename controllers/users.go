package controllers

import (
	"net/http"
	"zucora/backend/models"
	"zucora/backend/services"

	"github.com/gin-gonic/gin"
)


type UserController struct {
    userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
    return &UserController{userService: userService}
}
type req struct{
    Name string `json:"name"`
    Email string `json:"email"`
    Password string `json:"password"`  
}
func (c *UserController) RegisterUser(ctx *gin.Context) {
    var user models.User
    var request req
    if err := ctx.ShouldBindJSON(&request); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    user.Name=request.Name
    user.Email=request.Email
    user.Password=request.Password
    if err := c.userService.Register(&user); err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
        return
    }

    ctx.JSON(http.StatusCreated, user)
}

func (c *UserController) GetUserByID(ctx *gin.Context) {
    userID := ctx.Param("id")
    // Convert userID to uint and handle errors...

    user, err := c.userService.GetUserByID(userID)
    if err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    ctx.JSON(http.StatusOK, user)
}