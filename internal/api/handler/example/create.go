// Package example holds the HTTP handlers for the example reference domain.
package example

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/tclavelloux/promy-template-go/internal/api/request"
	"github.com/tclavelloux/promy-template-go/internal/domain/example"
	"github.com/tclavelloux/promy-template-go/internal/domain/example/model"
)

// CreateHandler handles POST /api/examples/v1/.
type CreateHandler struct {
	extractor request.Extractor
	service   example.Service
}

// NewCreateHandler creates a new CreateHandler.
func NewCreateHandler(extractor request.Extractor, service example.Service) *CreateHandler {
	return &CreateHandler{
		extractor: extractor,
		service:   service,
	}
}

// Handle creates a new example.
func (h *CreateHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var req model.CreateExampleRequest
		if err := h.extractor.Extract(c, &req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid request format")
		}

		created, err := h.service.Create(ctx, &req)
		if err != nil {
			return handleServiceError(err, "failed to create example")
		}

		return c.JSON(http.StatusCreated, created)
	}
}

// handleServiceError maps example service errors to appropriate HTTP responses.
func handleServiceError(err error, fallback string) error {
	errMsg := err.Error()

	switch {
	case strings.Contains(errMsg, "invalid example"):
		return echo.NewHTTPError(http.StatusBadRequest, errMsg)
	case strings.Contains(errMsg, "not found"):
		return echo.NewHTTPError(http.StatusNotFound, errMsg)
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, fallback)
	}
}
