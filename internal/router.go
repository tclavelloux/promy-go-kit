package internal

import (
	"github.com/ankorstore/yokai/fxhttpserver"
	exampleHandler "github.com/tclavelloux/promy-template-go/internal/api/handler/example"
	"github.com/tclavelloux/promy-template-go/internal/api/request"
	"go.uber.org/fx"
)

// Router is used to register the application HTTP middlewares and handlers.
func Router() fx.Option {
	return fx.Options(
		// shared HTTP request extraction, consumed by handlers
		fx.Provide(
			fx.Annotate(request.NewRequestParametersExtractor, fx.As(new(request.Extractor))),
		),

		// example domain: seeded reference pattern, delete once real domains exist
		fxhttpserver.AsHandlersGroup(
			"/api/examples",
			[]*fxhttpserver.HandlerRegistration{
				fxhttpserver.NewHandlerRegistration("POST", "/v1/", exampleHandler.NewCreateHandler),
				fxhttpserver.NewHandlerRegistration("GET", "/v1/:id", exampleHandler.NewGetHandler),
			},
		),
	)
}
