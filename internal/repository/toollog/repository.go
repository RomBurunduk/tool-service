package toollog

import (
	"context"
	"fmt"

	"tool-service/internal/model"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Insert(ctx context.Context, tc model.ToolCall) error {
	const q = `
		INSERT INTO tool_calls (request_id, tool_name, input, output, status, error, duration_ms)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.ExecContext(ctx, q,
		tc.RequestID,
		tc.ToolName,
		tc.Input,
		tc.Output,
		string(tc.Status),
		tc.Error,
		tc.DurationMs,
	)
	if err != nil {
		return fmt.Errorf("insert tool_call: %w", err)
	}
	return nil
}
