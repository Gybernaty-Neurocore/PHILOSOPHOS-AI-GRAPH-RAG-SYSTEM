package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func GetEmbedding(text string) ([]float32, error) {
	hfToken := os.Getenv("HF_TOKEN")
	if hfToken == "" {
		return nil, fmt.Errorf("HF_TOKEN not set")
	}

	payload := map[string]interface{}{
		"inputs": text,
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "https://api-inference.huggingface.co/models/Qwen/Qwen1.5-Chat", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+hfToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Qwen embedding request failed with status %s: %s", resp.Status, string(b))
	}

	var result struct {
		Embedding []float32 `json:"embedding"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return result.Embedding, nil
}
