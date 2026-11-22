package internal

import (
	"github.com/ankorstore/yokai/fxvalidator"
	"github.com/tclavelloux/promy-template-go/internal/domain/example"
	"go.uber.org/fx"
)

// Register is used to register the application dependencies.
func Register() fx.Option {
	return fx.Options(
		// validator module
		fxvalidator.FXValidatorModule,

		// example domain: seeded reference pattern, delete once real domains exist
		example.RegisterExampleComponents(),
	)
}
