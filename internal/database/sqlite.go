package database

import (
	"database/sql"
	"fmt"
	"sync/atomic"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/richardbizik/gommentary/internal/database/queries"
	_ "modernc.org/sqlite"
)

type Config struct {
	File            string `yaml:"file" json:"file" env:"DB_FILE" env-default:"gommentary.db"`
	MigrationsDir   string `yaml:"migrationsDir" json:"migrationsDir" env:"DB_MIGRATIONS_DIR" env-default:"./sql/migrations"`
	MigrationsTable string `yaml:"migrationTable" json:"migrationTable" env:"MIGRATION_TABLE" env-default:"schema_migrations"`
}

type Sqlite struct {
	db       *sql.DB
	Queries  *queries.Queries
	isLocked atomic.Bool
	config   *Config
}

func NewDatabase(config Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite", config.File)
	if err != nil {
		return nil, err
	}
	driver, err := sqlite.WithInstance(db, &sqlite.Config{
		MigrationsTable: config.MigrationsTable,
		DatabaseName:    "",
		NoTxWrap:        false,
	})
	if err != nil {
		return nil, err
	}
	migrator, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s", config.MigrationsDir), "", driver)
	if err != nil {
		return nil, err
	}
	migrator.Up()
	_, err = db.Exec("PRAGMA journal_mode=WAL;")
	if err != nil {
		return nil, err
	}

	queries := queries.New(db)
	return &Sqlite{
		db:       db,
		Queries:  queries,
		isLocked: atomic.Bool{},
		config:   &Config{},
	}, nil
}

func (db *Sqlite) Close() error {
	return db.db.Close()
}

func (db *Sqlite) Tx() (*sql.Tx, error) {
	return db.db.Begin()
}
