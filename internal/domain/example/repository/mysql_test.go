package repository_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tclavelloux/promy-template-go/internal"
	"github.com/tclavelloux/promy-template-go/internal/domain/example/model"
	"github.com/tclavelloux/promy-template-go/internal/domain/example/repository"
	"go.uber.org/fx"
)

func runRepositoryTest(t *testing.T) (repository.Repository, *sql.DB, context.Context) {
	t.Helper()

	var repo repository.Repository
	var db *sql.DB

	internal.RunTest(
		t,
		fx.Populate(&repo, &db),
	)

	return repo, db, t.Context()
}

func TestRepository_Create(t *testing.T) {
	t.Run("creates an example row", func(t *testing.T) {
		repo, db, ctx := runRepositoryTest(t)

		example := &model.Example{
			ID:        "test-example-001",
			Name:      "Created Example",
			CreatedAt: time.Now(),
		}

		err := repo.Create(ctx, example)
		require.NoError(t, err)

		var count int
		err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM examples WHERE id = ?", example.ID).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 1, count)
	})

	t.Run("round-trips through GetByID", func(t *testing.T) {
		repo, _, ctx := runRepositoryTest(t)

		example := &model.Example{
			ID:        "test-example-002",
			Name:      "Round Trip Example",
			CreatedAt: time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC),
		}

		require.NoError(t, repo.Create(ctx, example))

		retrieved, err := repo.GetByID(ctx, example.ID)
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.Equal(t, example.ID, retrieved.ID)
		assert.Equal(t, example.Name, retrieved.Name)
	})
}

func TestRepository_GetByID(t *testing.T) {
	t.Run("returns error for non-existent example", func(t *testing.T) {
		repo, _, ctx := runRepositoryTest(t)

		found, err := repo.GetByID(ctx, "non-existent-id")
		require.Error(t, err)
		assert.Nil(t, found)
		assert.Contains(t, err.Error(), "example with ID 'non-existent-id' not found")
	})
}
