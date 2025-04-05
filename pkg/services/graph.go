package services

import (
	"context"
	"fmt"
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
	// Для упрощения пока просто вернем context
	return fmt.Sprintf("Контекст:\n%s\n\nВопрос: %s", context, query), nil
}
