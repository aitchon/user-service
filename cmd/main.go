package main

import (
	"log"

	echoSwagger "github.com/swaggo/echo-swagger"
	"user-service/controllers"
	"user-service/db"
	_ "user-service/docs"
	"user-service/repositories"
	"user-service/services"

	"github.com/labstack/echo/v4"
)

func main() {
	// Initialize SQLite database connection
	database := db.ConnectDB()
	defer database.Close()

	// Set up repository and service
	userRepo := repositories.NewUserRepository(database)
	userService := services.NewUserService(userRepo)

	// Initialize Echo
	e := echo.New()

	// Routes
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/users", controllers.GetUsers(userService))
	e.POST("/users", controllers.CreateUser(userService))
	e.PUT("/users/:id", controllers.UpdateUser(userService))
	e.DELETE("/users/:id", controllers.DeleteUser(userService))

	// Start server
	log.Fatal(e.Start(":3002"))
}
