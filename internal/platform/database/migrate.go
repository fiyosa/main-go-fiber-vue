package database

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func MigrateUp(dsn string) error {
	if err := EnsureDatabaseExists(dsn); err != nil {
		return err
	}
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

// EnsureDatabaseExists creates the target database if it does not exist yet.
func EnsureDatabaseExists(dsn string) error {
	u, err := url.Parse(dsn)
	if err != nil {
		return fmt.Errorf("parse dsn: %w", err)
	}

	dbName := strings.TrimPrefix(u.Path, "/")
	if dbName == "" || dbName == "postgres" {
		return nil
	}

	u.Path = "/postgres"
	rootDSN := u.String()

	db, err := sql.Open("postgres", rootDSN)
	if err != nil {
		return fmt.Errorf("open root db: %w", err)
	}
	defer db.Close()

	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)"
	if err := db.QueryRow(query, dbName).Scan(&exists); err != nil {
		return fmt.Errorf("check database exists: %w", err)
	}

	if !exists {
		quotedDB := fmt.Sprintf("%q", dbName)
		if _, err := db.Exec("CREATE DATABASE " + quotedDB); err != nil {
			return fmt.Errorf("create database %s: %w", dbName, err)
		}
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
