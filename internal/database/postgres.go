// Package database - пакет в котором хранятся структуры, отвечающие за подключение к БД
package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/salex06/pr-service/internal/config"
)

// DB представляет собой структуру,
// хранящую пул соединений к БД PostgreSQL
type DB struct {
	Pool *pgxpool.Pool
}

// NewDB конструирует на основе конфига
// соединение к БД и возвращает объект DB
func NewDB(cfg *config.DBConfig) (*DB, error) {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("unable to parse connection string: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL")
	return &DB{Pool: pool}, nil
}

// Close закрывает соединение с БД
func (db *DB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}
