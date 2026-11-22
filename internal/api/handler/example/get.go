package example

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tclavelloux/promy-template-go/internal/domain/example"
)

// GetHandler handles GET /api/examples/v1/:id.
type GetHandler struct {
	service example.Service
}

// NewGetHandler creates a new GetHandler.
func NewGetHandler(service example.Service) *GetHandler {
	return &GetHandler{service: service}
}

// Handle returns an example by its ID.
func (h *GetHandler) Handle() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		id := c.Param("id")

		found, err := h.service.GetByID(ctx, id)
		if err != nil {
			return handleServiceError(err, "failed to get example")
		}

		return c.JSON(http.StatusOK, found)
	}
}
