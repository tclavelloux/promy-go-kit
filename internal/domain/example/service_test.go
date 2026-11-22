package example_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tclavelloux/promy-template-go/internal"
	"github.com/tclavelloux/promy-template-go/internal/domain/example"
	"github.com/tclavelloux/promy-template-go/internal/domain/example/model"
	"go.uber.org/fx"
)

func runServiceTest(t *testing.T) (example.Service, context.Context) {
	t.Helper()

	var service example.Service

	internal.RunTest(
		t,
		fx.Populate(&service),
	)

	return service, t.Context()
}

func TestService_Create(t *testing.T) {
	t.Run("assigns id and created_at, then persists", func(t *testing.T) {
		service, ctx := runServiceTest(t)

		created, err := service.Create(ctx, &model.CreateExampleRequest{Name: "First Example"})
		require.NoError(t, err)
		require.NotNil(t, created)
		assert.NotEmpty(t, created.ID)
		assert.Equal(t, "First Example", created.Name)
		assert.False(t, created.CreatedAt.IsZero())

		fetched, err := service.GetByID(ctx, created.ID)
		require.NoError(t, err)
		assert.Equal(t, created.ID, fetched.ID)
		assert.Equal(t, "First Example", fetched.Name)
	})

	t.Run("rejects a missing mandatory name", func(t *testing.T) {
		service, ctx := runServiceTest(t)

		created, err := service.Create(ctx, &model.CreateExampleRequest{Name: ""})
		require.Error(t, err)
		assert.Nil(t, created)
		assert.Contains(t, err.Error(), "invalid example")
	})
}

func TestService_GetByID(t *testing.T) {
	t.Run("returns not found for unknown id", func(t *testing.T) {
		service, ctx := runServiceTest(t)

		found, err := service.GetByID(ctx, "does-not-exist")
		require.Error(t, err)
		assert.Nil(t, found)
		assert.Contains(t, err.Error(), "not found")
	})
}
