package repository

import (
	"context"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func Connect(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("connect db: %w", err)
	}
	return db, nil
}

func Migrate(ctx context.Context, db *sqlx.DB) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS phones (
			id BIGSERIAL PRIMARY KEY,
			brand TEXT NOT NULL,
			model TEXT NOT NULL,
			price DOUBLE PRECISION,
			payload JSONB NOT NULL
		)`,
		`CREATE INDEX IF NOT EXISTS phones_brand_model_idx ON phones (brand, model)`,
		`CREATE EXTENSION IF NOT EXISTS pg_trgm`,
		`CREATE INDEX IF NOT EXISTS phones_fullname_trgm_idx ON phones USING gin ((lower(trim(brand) || ' ' || trim(model))) gin_trgm_ops)`,
		`CREATE TABLE IF NOT EXISTS tool_calls (
			id BIGSERIAL PRIMARY KEY,
			request_id UUID NOT NULL,
			tool_name TEXT NOT NULL,
			input JSONB,
			output JSONB,
			status TEXT NOT NULL,
			error TEXT,
			duration_ms INT,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,
	}
	for _, s := range stmts {
		if _, err := db.ExecContext(ctx, s); err != nil {
			return fmt.Errorf("migrate: %w", err)
		}
	}
	return nil
}
