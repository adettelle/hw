package migrator

import (
	"database/sql"
	"embed"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5" // import pgx/v5 driver for golang-migrate
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"go.uber.org/zap"
)

// migrationsDir specifies the directory containing migration files within the embedded filesystem.
const migrationsDir = "migration"

// MigrationFS is an embedded filesystem containing SQL migration files.
//
// The //go:embed directive embeds all SQL files in the `migration` directory into the binary.
//
//go:embed migration/*.sql
var MigrationsFS embed.FS

// MustApplyMigrations applies all pending migrations to the PostgreSQL database specified by dbParams.
// The migrations are sourced from the embedded `MigrationFS`.
//
// Parameters:
//   - dbParams: A connection string containing database configuration details.
func MustApplyMigrations(dbParams string, logg *zap.Logger) {
	// Create a new source driver from the embedded filesystem
	srcDriver, err := iofs.New(MigrationsFS, migrationsDir)
	if err != nil {
		logg.Fatal("unable to create a new source driver from the embedded filesystem", zap.Error(err))
	}
	// Open the database connection
	db, err := sql.Open("pgx", dbParams)
	if err != nil {
		logg.Fatal("unable to open the database connection", zap.Error(err))
	}

	// Create a PostgreSQL driver instance
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		logg.Fatal("unable to create db instance", zap.Error(err))
	}

	// Create a new migrator instance with the embedded migration files
	migrator, err := migrate.NewWithInstance("migration_embedded_sql_files", srcDriver, "psql_db", driver)
	if err != nil {
		logg.Fatal("unable to create migration", zap.Error(err))
	}

	// Apply all migrations; ignore the error if there are no changes
	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logg.Fatal("unable to apply migration", zap.Error(err))
	}

	migrator.Close()

	logg.Info("Migrations applied")
}
