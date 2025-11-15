package test_env

import (
	"context"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func NewTestDatabaseConnection(ctx context.Context) (*sqlx.DB, error) {

	cfg := struct {
		Host     string
		Port     int
		User     string
		Password string
		Database string
		SSLMode  string
	}{
		Host:     "localhost",
		Port:     5435,
		User:     "postgres",
		Password: "postgres",
		Database: "reviewers_test",
		SSLMode:  "disable",
	}

	dsn := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Database, cfg.User, cfg.Password, cfg.SSLMode,
	)

	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := db.PingContext(pingCtx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return db, nil
}
