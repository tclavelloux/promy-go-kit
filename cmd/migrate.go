package cmd

import (
	"strings"

	"github.com/ankorstore/yokai/fxcore"
	"github.com/ankorstore/yokai/fxsql"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

func init() {
	rootCmd.AddCommand(migrateDatabaseCmd)
}

var migrateDatabaseCmd = &cobra.Command{
	Use:     "migrate-database [command] [args...]",
	Short:   "Execute database migrations",
	Example: strings.Join(migrateExamples, "\n"),
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fxcore.
			NewBootstrapper().
			WithContext(cmd.Context()).
			WithOptions(migrateOptions(args)...).
			RunApp()
	},
}

// migrateOptions builds the fx options of migrate-database. The wait invoke
// must sit between the SQL module and the migration: root invokes run in order.
func migrateOptions(args []string) []fx.Option {
	opts := []fx.Option{
		fx.NopLogger,
		// modules
		fxsql.FxSQLModule,
	}

	if needsDatabaseWait(args[0]) {
		// wait for the database to accept connections
		opts = append(opts, fx.Invoke(waitForDatabase))
	}

	// migrate and shutdown
	return append(opts, fxsql.RunFxSQLMigrationAndShutdown(args[0], args[1:]...))
}

var migrateExamples = []string{
	"  migrate-database up                    migrate the DB to the most recent version available",
	"  migrate-database up-by-one             migrate the DB up by 1",
	"  migrate-database up-to VERSION         migrate the DB to a specific VERSION",
	"  migrate-database down                  roll back the version by 1",
	"  migrate-database down-to VERSION       roll back to a specific VERSION",
	"  migrate-database redo                  re-run the latest migration",
	"  migrate-database reset                 roll back all migrations",
	"  migrate-database status                dump the migration status for the current DB",
	"  migrate-database version               print the current version of the database",
	"  migrate-database create NAME [sql|go]  creates new migration file with the current timestamp",
	"  migrate-database fix                   apply sequential ordering to migrations",
	"  migrate-database validate              check migration files without running them",
}
