package phone

import (
	"context"
	"database/sql"
	"errors"
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

func (r *Repository) GetByBrandModel(ctx context.Context, brand, modelName string) (*model.Phone, error) {
	const q = `
		SELECT id, brand, model, price, payload
		FROM phones
		WHERE brand = $1 AND model = $2
		ORDER BY id
		LIMIT 1
	`
	var row struct {
		ID      int64           `db:"id"`
		Brand   string          `db:"brand"`
		Model   string          `db:"model"`
		Price   sql.NullFloat64 `db:"price"`
		Payload []byte          `db:"payload"`
	}
	err := r.db.GetContext(ctx, &row, q, brand, modelName)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get phone: %w", err)
	}
	p := &model.Phone{
		ID:      row.ID,
		Brand:   row.Brand,
		Model:   row.Model,
		Payload: row.Payload,
	}
	if row.Price.Valid {
		v := row.Price.Float64
		p.Price = &v
	}
	return p, nil
}

func (r *Repository) Truncate(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `TRUNCATE phones`)
	return err
}

func (r *Repository) InsertBatch(ctx context.Context, rows []model.Phone) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	const q = `INSERT INTO phones (brand, model, price, payload) VALUES ($1, $2, $3, $4)`
	for _, p := range rows {
		if _, err := tx.ExecContext(ctx, q, p.Brand, p.Model, p.Price, p.Payload); err != nil {
			return fmt.Errorf("insert phone: %w", err)
		}
	}
	return tx.Commit()
}

// InsertOne for tests if needed
func (r *Repository) InsertOne(ctx context.Context, p model.Phone) error {
	payload := p.Payload
	if len(payload) == 0 {
		payload = []byte("{}")
	}
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO phones (brand, model, price, payload) VALUES ($1, $2, $3, $4)`,
		p.Brand, p.Model, p.Price, payload,
	)
	return err
}
