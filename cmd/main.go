package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Mukam21/RAG_server-Golang/pkg/handlers"
	"github.com/Mukam21/RAG_server-Golang/pkg/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables:", err)
	} else {
		log.Println("Successfully loaded .env file")
	}

	if os.Getenv("PG_PASSWORD") == "" {
		log.Fatal("PG_PASSWORD is not set")
	}

	if err := services.InitDB(); err != nil {
		log.Fatal("Failed to init DB:", err)
	}
	defer services.CloseConnection()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1"})
	r.POST("/upload", handlers.UploadDocument)
	r.POST("/add", handlers.AddDocuments)
	r.POST("/query", handlers.Query)

	fmt.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
