package main

import (
	"flag"
	"fmt"
	"log"
	"zucora/backend/controllers"
	"zucora/backend/database"
	"zucora/backend/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	syncdb := flag.Bool("syncdb", false, "sync database boolean param")
	flag.Parse()

	database.ConnectDB()

	if *syncdb {
		database.SyncDb()
	}

}

func main() {
	fmt.Println("Hello, world.")
	router := gin.Default()

	router.POST("/create/admin", controllers.CreateSuperAdmin)
	router.POST("/admin/login", controllers.SuperAdminLogin)
	router.POST("/user/login", controllers.UserLogin)

	//Common API functions
	router.PUT("/changedpwd", middleware.AuthRequired, controllers.ChangePassword)
	router.GET("/logout", controllers.Logout)
	// Group for superadmin routes
	adminGroup := router.Group("/admin")
	adminGroup.Use(middleware.AuthRequired) // Custom middleware to check admin role
	{
		adminGroup.POST("/create/user", controllers.CreateUser)
		adminGroup.GET("/users", controllers.GetAllUsers)
		adminGroup.GET("/user/:id", controllers.GetUserDetail)
		adminGroup.DELETE("/user/:id", controllers.DeleteUser)
	}
	// Group for user routes
	userGroup := router.Group("/user")
	userGroup.Use(middleware.AuthRequired) // Custom middleware to check user role
	{
		userGroup.POST("/verify2fa", controllers.VerifyTwoFA)
		userGroup.GET("/:id", controllers.GetUserDetail)
	}
	router.Run()
}
