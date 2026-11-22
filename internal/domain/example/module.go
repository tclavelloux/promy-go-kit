package example

import (
	"github.com/tclavelloux/promy-template-go/internal/domain/example/repository"
	"go.uber.org/fx"
)

// RegisterExampleComponents registers the example domain's repository and service.
func RegisterExampleComponents() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotate(repository.NewRepository, fx.As(new(repository.Repository))),
		),
		fx.Provide(
			fx.Annotate(NewService, fx.As(new(Service))),
		),
	)
}
