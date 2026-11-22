package example_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tclavelloux/promy-template-go/internal"
	"github.com/tclavelloux/promy-template-go/internal/domain/example/model"
	"go.uber.org/fx"
)

func setupServer(tb testing.TB) (*echo.Echo, *sql.DB) {
	tb.Helper()

	var httpServer *echo.Echo
	var db *sql.DB

	internal.RunTest(
		tb,
		fx.Populate(&httpServer, &db),
	)

	return httpServer, db
}

func TestCreateHandler(t *testing.T) {
	httpServer, db := setupServer(t)

	t.Run("success", func(t *testing.T) {
		body, err := json.Marshal(model.CreateExampleRequest{Name: "New Example"})
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/examples/v1/", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		httpServer.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var created model.Example
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &created))
		assert.NotEmpty(t, created.ID)
		assert.Equal(t, "New Example", created.Name)
		assert.False(t, created.CreatedAt.IsZero())

		// DB side-effect: the row was actually persisted.
		var count int
		require.NoError(t, db.QueryRow("SELECT COUNT(*) FROM examples WHERE id = ?", created.ID).Scan(&count))
		assert.Equal(t, 1, count)
	})

	t.Run("missing mandatory name returns 400", func(t *testing.T) {
		body, err := json.Marshal(model.CreateExampleRequest{Name: ""})
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/examples/v1/", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		httpServer.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "invalid example")
	})
}

func TestGetHandler(t *testing.T) {
	httpServer, _ := setupServer(t)

	t.Run("round trip: create then get", func(t *testing.T) {
		body, err := json.Marshal(model.CreateExampleRequest{Name: "Round Trip"})
		require.NoError(t, err)

		createReq := httptest.NewRequest(http.MethodPost, "/api/examples/v1/", bytes.NewBuffer(body))
		createReq.Header.Set("Content-Type", "application/json")
		createRec := httptest.NewRecorder()
		httpServer.ServeHTTP(createRec, createReq)
		require.Equal(t, http.StatusCreated, createRec.Code)

		var created model.Example
		require.NoError(t, json.Unmarshal(createRec.Body.Bytes(), &created))

		getReq := httptest.NewRequest(http.MethodGet, "/api/examples/v1/"+created.ID, nil)
		getRec := httptest.NewRecorder()
		httpServer.ServeHTTP(getRec, getReq)

		assert.Equal(t, http.StatusOK, getRec.Code)

		var fetched model.Example
		require.NoError(t, json.Unmarshal(getRec.Body.Bytes(), &fetched))
		assert.Equal(t, created.ID, fetched.ID)
		assert.Equal(t, "Round Trip", fetched.Name)
	})

	t.Run("unknown id returns 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/examples/v1/does-not-exist", nil)
		rec := httptest.NewRecorder()
		httpServer.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}
