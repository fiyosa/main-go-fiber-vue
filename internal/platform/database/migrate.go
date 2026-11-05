package database

import (
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"runtime"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func MigrateUp(dsn string) error {
	m, err := newMigrator(dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}

func MigrateDown(dsn string, steps int) error {
	if steps <= 0 {
		steps = 1
	}
	m, err := newMigrator(dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Steps(-steps); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}

func MigrateStatus(dsn string) (uint, bool, error) {
	m, err := newMigrator(dsn)
	if err != nil {
		return 0, false, err
	}
	defer m.Close()

	version, dirty, err := m.Version()
	if errors.Is(err, migrate.ErrNilVersion) {
		return 0, false, nil
	}
	return version, dirty, err
}

func newMigrator(dsn string) (*migrate.Migrate, error) {
	sourceURL, err := migrationSourceURL()
	if err != nil {
		return nil, err
	}
	m, err := migrate.New(sourceURL, dsn)
	if err != nil {
		return nil, fmt.Errorf("open migrator: %w", err)
	}
	return m, nil
}

func migrationSourceURL() (string, error) {
	dir, err := filepath.Abs("migrations")
	if err != nil {
		return "", fmt.Errorf("resolve migrations dir: %w", err)
	}
	path := filepath.ToSlash(dir)
	if runtime.GOOS == "windows" {
		// The migrate file driver joins the URL host and path, so a drive
		// letter must be the host: file://C:/project/migrations.
		return "file://" + path, nil
	}
	return (&url.URL{Scheme: "file", Path: path}).String(), nil
}
