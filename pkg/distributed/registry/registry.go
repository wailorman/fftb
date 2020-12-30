package registry

import (
	"github.com/pkg/errors"

	// "database/sql"

	// "gorm.io/driver/sqlite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	// "github.com/golang-migrate/migrate"
	// "github.com/golang-migrate/migrate/database/sqlite3"

	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"

	// ss
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"

	// "database/sql"
	// _ "github.com/mattn/go-sqlite3"
	// "github.com/golang-migrate/migrate/v4"
	// "github.com/golang-migrate/migrate/v4/database/sqlite3"
	// _ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/subchen/go-trylock/v2"
)

// SqliteRegistry _
type SqliteRegistry struct {
	db              *sql.DB
	gdb             *gorm.DB
	freeSegmentLock trylock.TryLocker
}

// NewSqliteRegistry _
func NewSqliteRegistry(databasePath, migrationsPath string) (*SqliteRegistry, error) {
	var err error

	// TODO: migrations
	// pkg/distributed/registry/migrations

	r := &SqliteRegistry{
		freeSegmentLock: trylock.New(),
	}

	r.db, err = sql.Open("sqlite3", databasePath)

	if err != nil {
		return nil, errors.Wrap(err, "Initializing sqlite database file")
	}

	// -
	driver, err := sqlite3.WithInstance(r.db, &sqlite3.Config{})

	if err != nil {
		return nil, errors.Wrap(err, "Initializing sqlite migrations driver")
	}

	// -
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"sqlite3",
		driver,
	)

	if err != nil {
		return nil, errors.Wrap(err, "Initializing migrator")
	}

	// -
	err = m.Steps(1)

	// -
	// if err != nil {
	// 	return nil, errors.Wrap(err, "Performing migrations")
	// }

	r.gdb, err = gorm.Open(sqlite.Open(databasePath), &gorm.Config{})

	if err != nil {
		return nil, errors.Wrap(err, "Initializing sqlite gorm db")
	}

	return r, nil
}
