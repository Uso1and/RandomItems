package main

import (
	"RandomItems/internal/app/handlers"
	"RandomItems/internal/domain/infrastructure/database"
	"RandomItems/internal/domain/repositories"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := database.Init(); err != nil {
		log.Fatalf("Error DB open!: %v", err)
	}
	defer database.DB.Close()

	log.Println("Database connected!")

	userRepo := repositories.NewUserRepository(database.DB)

	userHandler := handlers.NewUserHandler(userRepo)

	r := gin.Default()

	r.POST("/user", userHandler.CreateUser)
	r.GET("/user/:id", userHandler.GetUser)
	r.Run(":8080")
}
