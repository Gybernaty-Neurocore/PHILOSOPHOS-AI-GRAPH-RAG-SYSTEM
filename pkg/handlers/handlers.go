package handlers

import (
	"io"
	"net/http"
	"strings"
	"sync"

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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var wg sync.WaitGroup
	for _, doc := range req.Documents {
		text := strings.TrimSpace(doc.Text)
		if len(text) < 5 {
			continue
		}
		wg.Add(1)
		go func(t string) {
			defer wg.Done()
			_ = services.AddGraphDocument(t)
		}(text)
	}
	wg.Wait()
	c.JSON(http.StatusOK, gin.H{"message": "Documents added to graph"})
}

type QueryRequest struct {
	Query string `json:"query" binding:"required"`
}

func Query(c *gin.Context) {
	var req QueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, err := services.QueryGraphRAG(req.Query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"response": response})
}

func UploadDocument(c *gin.Context) {
	file, _, err := c.Request.FormFile("document")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file read error"})
		return
	}
	defer file.Close()
	content, _ := io.ReadAll(file)
	_ = services.AddGraphDocument(string(content))
	c.JSON(http.StatusOK, gin.H{"message": "File uploaded to graph"})
}
