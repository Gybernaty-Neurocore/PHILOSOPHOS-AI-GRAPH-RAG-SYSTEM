package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func AddGraphDocument(content string) error {
	embedding, err := GetEmbedding(content)
	if err != nil {
		return err
	}
	var id int
	err = pgConn.QueryRow(context.Background(), `INSERT INTO nodes (content, embedding) VALUES ($1, $2) RETURNING id`, content, embedding).Scan(&id)
	if err != nil {
		return err
	}
	rows, _ := pgConn.Query(context.Background(), `SELECT id FROM nodes WHERE id != $1 LIMIT 1`, id)
	for rows.Next() {
		var otherID int
		_ = rows.Scan(&otherID)
		_, _ = pgConn.Exec(context.Background(), `INSERT INTO edges (from_node, to_node, relation) VALUES ($1, $2, $3)`, id, otherID, "related")
	}
	return nil
}

func QueryGraphRAG(query string) (string, error) {
	embedding, err := GetEmbedding(query)
	if err != nil {
		return "", err
	}
	var content string
	var nodeID int
	err = pgConn.QueryRow(context.Background(), `SELECT id, content FROM nodes ORDER BY embedding <-> $1 LIMIT 1`, embedding).Scan(&nodeID, &content)
	if err != nil {
		return "", err
	}
	rows, _ := pgConn.Query(context.Background(), `SELECT content FROM nodes WHERE id IN (SELECT to_node FROM edges WHERE from_node = $1)`, nodeID)
	var neighborContents []string
	for rows.Next() {
		var n string
		_ = rows.Scan(&n)
		neighborContents = append(neighborContents, n)
	}
	contextText := content + "\n" + strings.Join(neighborContents, "\n")
	return GenerateResponse(query, contextText)
}

func GenerateResponse(query, context string) (string, error) {
	hfToken := os.Getenv("HF_TOKEN")
	if hfToken == "" {
		return "", fmt.Errorf("HF_TOKEN not set")
	}

	payload := map[string]interface{}{
		"inputs": map[string]string{
			"prompt": fmt.Sprintf("Вот контекст:\n%s\n\nОтветь на вопрос: %s", context, query),
		},
	}

	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "https://api-inference.huggingface.co/models/Qwen/Qwen1.5-Chat", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+hfToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("HF generation failed with status %s: %s", resp.Status, string(b))
	}

	var result []struct {
		GeneratedText string `json:"generated_text"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	if len(result) > 0 {
		return result[0].GeneratedText, nil
	}

	return "", fmt.Errorf("empty response from LLM")
}
