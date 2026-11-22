package internal

import (
	"context"
	"fmt"
	"testing"

	"github.com/ankorstore/yokai-contrib/fxgomysqlserver"
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxcore"
	"github.com/ankorstore/yokai/fxhttpserver"
	"github.com/ankorstore/yokai/fxsql"
	"go.uber.org/fx"
)

func init() {
	RootDir = fxcore.RootDir(1)
}

// RootDir is the application root directory.
var RootDir string

// Bootstrapper can be used to load modules, options, dependencies, routing and bootstraps the application.
var Bootstrapper = fxcore.NewBootstrapper().WithOptions(
	// yokai modules registration
	fxsql.FxSQLModule,
	fxhttpserver.FxHttpServerModule,

	// dependencies registration
	Register(),
	// routing registration
	Router(),
)

// Run starts the application, with a provided [context.Context].
func Run(ctx context.Context) {
	Bootstrapper.WithContext(ctx).RunApp()
}

// RunTest starts the application in test mode, with an optional list of [fx.Option].
func RunTest(tb testing.TB, options ...fx.Option) {
	tb.Helper()

	tb.Setenv("MODULES_SQL_MIGRATIONS_PATH", fmt.Sprintf("%s/db/migrations", RootDir))

	Bootstrapper.RunTestApp(
		tb,
		// config lookup
		fxconfig.AsConfigPath(fmt.Sprintf("%s/configs/", RootDir)),
		// in-memory MySQL for tests (no docker/real MySQL required)
		fxgomysqlserver.FxGoMySQLServerModule,
		// run database migrations
		fxsql.RunFxSQLMigration("up"),
		// test options
		fx.Options(options...),
	)
}
