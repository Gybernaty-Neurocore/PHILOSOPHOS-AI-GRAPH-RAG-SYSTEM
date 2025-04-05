package main

import (
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

	if os.Getenv("HF_TOKEN") == "" {
		log.Fatal("HF_TOKEN is not set")
	}
	if err := services.InitDB(); err != nil {
		log.Fatal("Failed to init DB:", err)
	}
	defer services.CloseConnection()

	r := gin.Default()
	r.POST("/upload", handlers.UploadDocument)
	r.POST("/add", handlers.AddDocuments)
	r.POST("/query", handlers.Query)

	log.Println("Server started on :8080")
	_ = r.Run(":8080")
}
