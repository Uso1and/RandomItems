package main

import (
	"RandomItems/internal/app/handlers"
	"RandomItems/internal/app/services"
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
	itemRepo := repositories.NewItemRepository(database.DB)
	dropRepo := repositories.NewDropRepository(database.DB)

	dropService := services.NewDropService(itemRepo, dropRepo, userRepo)

	userHandler := handlers.NewUserHandler(userRepo)
	dropHandler := handlers.NewDropHandler(dropService)

	r := gin.Default()

	r.POST("/user", userHandler.CreateUser)
	r.GET("/user/:id", userHandler.GetUser)

	r.POST("/drop/:user_id", dropHandler.GenerateDrop)
	r.Run(":8080")
}
