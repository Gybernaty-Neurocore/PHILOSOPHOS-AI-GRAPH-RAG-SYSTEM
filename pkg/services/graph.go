package services

import (
	"context"
	"fmt"
	"math"
	"strings"
)

func AddGraphDocument(content string) error {
	fmt.Printf("Adding document with content: %q\n", content)
	embedding, err := GetEmbedding(content)
	if err != nil {
		return fmt.Errorf("failed to get embedding: %v", err)
	}

	var nodeID int
	err = pgConn.QueryRow(context.Background(), `
        INSERT INTO nodes (content, embedding) VALUES ($1, $2) RETURNING id
    `, content, fmt.Sprintf("[%s]", joinFloats(embedding))).Scan(&nodeID)
	if err != nil {
		return fmt.Errorf("failed to insert node: %v", err)
	}
	fmt.Printf("Node inserted successfully, nodeID: %d\n", nodeID)

	fmt.Printf("Document added successfully, nodeID: %d\n", nodeID)
	return nil
}

func QueryGraphRAG(query string) (string, error) {
	fmt.Printf("Starting QueryGraphRAG with query: %q\n", query)
	embedding, err := GetEmbedding(query)
	if err != nil {
		return "", fmt.Errorf("failed to get embedding: %v", err)
	}

	rows, err := pgConn.Query(context.Background(), `
        SELECT content, embedding FROM nodes
    `)
	if err != nil {
		return "", fmt.Errorf("failed to query nodes: %v", err)
	}
	defer rows.Close()

	var bestContent string
	maxSimilarity := -1.0

	for rows.Next() {
		var content, embeddingStr string
		if err := rows.Scan(&content, &embeddingStr); err != nil {
			return "", fmt.Errorf("failed to scan node: %v", err)
		}

		nodeEmbedding, err := parseEmbedding(embeddingStr)
		if err != nil {
			return "", fmt.Errorf("failed to parse embedding: %v", err)
		}

		similarity := cosineSimilarity(embedding, nodeEmbedding)
		fmt.Printf("Similarity for content %q: %f\n", content, similarity)
		if similarity > maxSimilarity {
			maxSimilarity = similarity
			bestContent = content
		}
	}
	if err := rows.Err(); err != nil {
		return "", fmt.Errorf("error iterating nodes: %v", err)
	}

	if bestContent == "" {
		return "", fmt.Errorf("no relevant content found")
	}

	fmt.Printf("Query result: %q\n", bestContent)
	return bestContent, nil
}

func joinFloats(floats []float32) string {
	parts := make([]string, len(floats))
	for i, f := range floats {
		parts[i] = fmt.Sprintf("%f", f)
	}
	return strings.Join(parts, ",")
}

func parseEmbedding(s string) ([]float32, error) {
	s = strings.Trim(s, "[]")
	parts := strings.Split(s, ",")
	result := make([]float32, len(parts))
	for i, p := range parts {
		var f float64
		if _, err := fmt.Sscanf(p, "%f", &f); err != nil {
			return nil, fmt.Errorf("failed to parse float at index %d: %v", i, err)
		}
		result[i] = float32(f)
	}
	return result, nil
}

func cosineSimilarity(a, b []float32) float64 {
	if len(a) != len(b) {
		return 0
	}
	var dotProduct, normA, normB float64
	for i := range a {
		dotProduct += float64(a[i]) * float64(b[i])
		normA += float64(a[i]) * float64(a[i])
		normB += float64(b[i]) * float64(b[i])
	}
	normA = math.Sqrt(normA)
	normB = math.Sqrt(normB)
	if normA == 0 || normB == 0 {
		return 0
	}
	return dotProduct / (normA * normB)
}
