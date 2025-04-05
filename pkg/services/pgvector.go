package services

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

var pgConn *pgx.Conn

// InitDB устанавливает соединение с PostgreSQL и создаёт таблицы для узлов и связей
func InitDB() error {
	connStr := fmt.Sprintf("postgres://postgres:%s@localhost:5438/postgres?sslmode=disable", os.Getenv("PG_PASSWORD"))
	var err error
	pgConn, err = pgx.Connect(context.Background(), connStr)
	if err != nil {
		return err
	}
	_, _ = pgConn.Exec(context.Background(), `CREATE EXTENSION IF NOT EXISTS vector`)
	_, _ = pgConn.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS nodes (id SERIAL PRIMARY KEY, content TEXT, embedding VECTOR(768))`)
	_, _ = pgConn.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS edges (id SERIAL PRIMARY KEY, from_node INT, to_node INT, relation TEXT)`)
	return nil
}

// CloseConnection закрывает соединение с базой данных
func CloseConnection() {
	_ = pgConn.Close(context.Background())
}
