package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Mukam21/RAG_server-Golang/pkg/services"
	"github.com/gin-gonic/gin"
)

type AddRequest struct {
	Documents []struct {
		Text string `json:"text" binding:"required"`
	} `json:"documents" binding:"required"`
}

func AddDocuments(c *gin.Context) {
	var req AddRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("Invalid JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid JSON: %v", err)})
		return
	}

	var errors []string

	for i, doc := range req.Documents {
		text := strings.TrimSpace(doc.Text)
		if len(text) < 5 {
			fmt.Printf("Skipping short document at index %d: %q\n", i, text)
			continue
		}
		fmt.Printf("Adding document %d: %q\n", i, text)
		if err := services.AddGraphDocument(text); err != nil {
			errors = append(errors, fmt.Sprintf("failed to add document %d '%s': %v", i, text, err))
		} else {
			fmt.Printf("Successfully added document %d: %q\n", i, text)
		}
	}

	if len(errors) > 0 {
		errMsg := strings.Join(errors, "; ")
		fmt.Printf("Errors in AddDocuments: %s\n", errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errMsg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Documents added to graph"})
}

type QueryRequest struct {
	Query string `json:"query" binding:"required"`
}

func Query(c *gin.Context) {
	var req QueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("Invalid JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid JSON: %v", err)})
		return
	}
	response, err := services.QueryGraphRAG(req.Query)
	if err != nil {
		fmt.Printf("Query failed: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("query failed: %v", err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"response": response})
}

func UploadDocument(c *gin.Context) {
	fmt.Println("Received upload request")

	file, header, err := c.Request.FormFile("document")
	if err != nil {
		fmt.Printf("Failed to get file: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("failed to get file: %v", err)})
		return
	}
	defer file.Close()

	fmt.Printf("Received file: %s, size: %d\n", header.Filename, header.Size)

	content, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("Failed to read file content: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("failed to read file content: %v", err)})
		return
	}

	fmt.Printf("File content: %q\n", string(content))

	sentences := strings.Split(strings.ReplaceAll(string(content), "\r\n", " "), ". ")
	var errors []string

	for i, sentence := range sentences {
		sentence = strings.TrimSpace(sentence)
		if len(sentence) < 5 {
			fmt.Printf("Skipping short sentence at index %d: %q\n", i, sentence)
			continue
		}
		fmt.Printf("Processing sentence %d: %q\n", i, sentence)
		if err := services.AddGraphDocument(sentence); err != nil {
			errors = append(errors, fmt.Sprintf("failed to add sentence %d '%s': %v", i, sentence, err))
		} else {
			fmt.Printf("Successfully added sentence %d: %q\n", i, sentence)
		}
	}

	if len(errors) > 0 {
		errMsg := strings.Join(errors, "; ")
		fmt.Printf("Errors in UploadDocument: %s\n", errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errMsg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded to graph"})
}
