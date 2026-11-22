// Package repository holds the data access layer for the example domain.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/huandu/go-sqlbuilder"
	"github.com/tclavelloux/promy-template-go/internal/domain/example/model"
)

// Repository is the data access contract for the example domain.
type Repository interface {
	Create(ctx context.Context, example *model.Example) error
	GetByID(ctx context.Context, id string) (*model.Example, error)
}

// DefaultRepository is the MySQL-backed implementation of Repository.
type DefaultRepository struct {
	db *sql.DB
}

// NewRepository creates a new DefaultRepository.
func NewRepository(db *sql.DB) *DefaultRepository {
	return &DefaultRepository{db: db}
}

var _ Repository = (*DefaultRepository)(nil)

// exampleColumns is the ordered column list shared by the read/write queries.
var exampleColumns = []string{"id", "name", "created_at"}

// Create persists a new example row.
func (r *DefaultRepository) Create(ctx context.Context, example *model.Example) error {
	ib := sqlbuilder.NewInsertBuilder()
	ib.InsertInto("examples")
	ib.Cols(exampleColumns...)
	ib.Values(example.ID, example.Name, example.CreatedAt)

	query, args := ib.Build()
	if _, err := r.db.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

// GetByID returns an example by its identifier.
func (r *DefaultRepository) GetByID(ctx context.Context, id string) (*model.Example, error) {
	sb := sqlbuilder.NewSelectBuilder()
	sb.Select(exampleColumns...)
	sb.From("examples")
	sb.Where(sb.Equal("id", id))

	query, args := sb.Build()
	row := r.db.QueryRowContext(ctx, query, args...)

	var example model.Example
	if err := row.Scan(&example.ID, &example.Name, &example.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("example with ID '%s' not found", id)
		}

		return nil, err
	}

	return &example, nil
}
