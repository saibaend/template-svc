package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/saibaend/template-svc/internal/app/someModule/model"
)

type ItemRepository struct {
	db *sqlx.DB
}

func NewItemRepository(db *sqlx.DB) *ItemRepository {
	return &ItemRepository{db: db}
}

func (r *ItemRepository) Create(ctx context.Context, item *model.Item) error {
	const query = `
		INSERT INTO items (title, description)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRowxContext(ctx, query, item.Title, item.Description).
		Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
}

func (r *ItemRepository) GetByID(ctx context.Context, id int64) (*model.Item, error) {
	const query = `
		SELECT id, title, description, created_at, updated_at
		FROM items
		WHERE id = $1`

	var item model.Item
	if err := r.db.GetContext(ctx, &item, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNotFound
		}

		return nil, fmt.Errorf("get item by id: %w", err)
	}

	return &item, nil
}

func (r *ItemRepository) List(ctx context.Context, limit, offset int) ([]model.Item, error) {
	const query = `
		SELECT id, title, description, created_at, updated_at
		FROM items
		ORDER BY id
		LIMIT $1 OFFSET $2`

	items := make([]model.Item, 0)
	if err := r.db.SelectContext(ctx, &items, query, limit, offset); err != nil {
		return nil, fmt.Errorf("list items: %w", err)
	}

	return items, nil
}

func (r *ItemRepository) Update(ctx context.Context, item *model.Item) error {
	const query = `
		UPDATE items
		SET title = $1, description = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING updated_at`

	if err := r.db.QueryRowxContext(ctx, query, item.Title, item.Description, item.ID).
		Scan(&item.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.ErrNotFound
		}

		return fmt.Errorf("update item: %w", err)
	}

	return nil
}

func (r *ItemRepository) Delete(ctx context.Context, id int64) (bool, error) {
	const query = `DELETE FROM items WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return false, fmt.Errorf("delete item: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("delete item rows affected: %w", err)
	}

	return rows > 0, nil
}
