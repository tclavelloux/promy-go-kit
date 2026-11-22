// Package request provides HTTP request-body extraction shared by handlers.
package request

import (
	"github.com/labstack/echo/v4"
)

// Extractor binds an HTTP request body onto params.
type Extractor interface {
	Extract(c echo.Context, params any) error
}

// RequestParametersExtractor is the default Extractor implementation.
type RequestParametersExtractor struct{}

// NewRequestParametersExtractor creates a new RequestParametersExtractor.
func NewRequestParametersExtractor() *RequestParametersExtractor {
	return &RequestParametersExtractor{}
}

var _ Extractor = (*RequestParametersExtractor)(nil)

// Extract binds the HTTP request body onto params. Validation happens in the
// service layer.
func (e *RequestParametersExtractor) Extract(c echo.Context, params any) error {
	return c.Bind(params)
}
