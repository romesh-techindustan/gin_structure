package main

import (
	"flag"
	"fmt"
	"log"
	"zucora/backend/controllers"
	"zucora/backend/database"
	"zucora/backend/repository"
	"zucora/backend/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	syncdb := flag.Bool("syncdb", true, "sync database boolean param")
	flag.Parse()

	database.ConnectDB()

	if *syncdb {
		database.SyncDb()
	}

}

func main() {
	fmt.Println("Hello, world.")
	router := gin.Default()
	db := database.ConnectDB()
	userRepository := repository.NewUserRepository(db)

	// Initialize services
	userService := services.NewUserService(userRepository)

	// Initialize controllers
	userController := controllers.NewUserController(userService)

	router.POST("/create", userController.RegisterUser)
	router.GET("/create/:id", userController.GetUserByID)

	router.Run()
}
