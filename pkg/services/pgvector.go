package services

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

var pgConn *pgx.Conn

func InitDB() error {
	connStr := fmt.Sprintf("postgres://postgres:%s@localhost:5438/postgres?sslmode=disable", os.Getenv("PG_PASSWORD"))
	var err error
	pgConn, err = pgx.Connect(context.Background(), connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	_, err = pgConn.Exec(context.Background(), `
        CREATE EXTENSION IF NOT EXISTS vector;
        CREATE TABLE IF NOT EXISTS nodes (
            id SERIAL PRIMARY KEY,
            content TEXT,
            embedding VECTOR(768)
        );
        CREATE TABLE IF NOT EXISTS edges (
            id SERIAL PRIMARY KEY,
            source_id INT REFERENCES nodes(id),
            target_id INT REFERENCES nodes(id),
            weight FLOAT
        );
    `)
	if err != nil {
		return fmt.Errorf("failed to create tables: %v", err)
	}

	fmt.Println("Database initialized successfully")
	return nil
}

func CloseConnection() {
	if pgConn != nil {
		pgConn.Close(context.Background())
	}
}
