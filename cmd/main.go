package main

import (
	"RandomItems/internal/domain/infrastructure/database"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := database.Init(); err != nil {
		log.Fatalf("Error DB open!: %v", err)
	}
	defer database.DB.Close()

	log.Println("Database connected!")

	r := gin.Default()

	r.Run(":8080")
}
