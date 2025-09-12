package migrator

import (
	"database/sql"
	"embed"
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
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
func MustApplyMigrations(dbParams string) {
	// Create a new source driver from the embedded filesystem
	srcDriver, err := iofs.New(MigrationsFS, migrationsDir)
	if err != nil {
		log.Fatal(err) // TODO HELP здесь должен быть zap Logger или нет ???
	}

	// Open the database connection
	db, err := sql.Open("pgx", dbParams)
	if err != nil {
		log.Fatal(err)
	}

	// Create a PostgreSQL driver instance
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("unable to create db instance: %v", err)
	}

	// Create a new migrator instance with the embedded migration files
	migrator, err := migrate.NewWithInstance("migration_embedded_sql_files", srcDriver, "psql_db", driver)
	if err != nil {
		log.Fatalf("unable to create migration: %v", err)
	}

	// Apply all migrations; ignore the error if there are no changes
	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("unable to apply migrations %v", err)
	}

	migrator.Close()

	log.Println("Migrations applied")
}
