package config

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	AppEnv         string `envconfig:"APP_ENV" default:"development"`
	AppPort        string `envconfig:"APP_PORT" default:"8080"`
	AppLocale      string `envconfig:"APP_LOCALE" default:"en"`
	AppMaintenance bool   `envconfig:"APP_MAINTENANCE" default:"false"`
	PublicDir      string `envconfig:"APP_PUBLIC_DIR" default:"./public"`

	DBHost     string `envconfig:"DB_HOST" default:"localhost"`
	DBPort     int    `envconfig:"DB_PORT" default:"5432"`
	DBDatabase string `envconfig:"DB_DATABASE" default:"go_fiber_vue"`
	DBUsername string `envconfig:"DB_USERNAME" default:"postgres"`
	DBPassword string `envconfig:"DB_PASSWORD"`
	DBSSLMode  string `envconfig:"DB_SSLMODE" default:"disable"`
	// DBLogSQL controls GORM SQL logging: true logs every executed query.
	DBLogSQL bool `envconfig:"DB_LOG_SQL" default:"false"`

	JWTSecret     string        `envconfig:"JWT_SECRET"`
	JWTAccessTTL  time.Duration `envconfig:"JWT_ACCESS_TTL" default:"3h"`
	JWTRefreshTTL time.Duration `envconfig:"JWT_REFRESH_TTL" default:"720h"`
}

func Load() (*Config, error) {
	loadDotEnv()

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}
	cfg.AppLocale = strings.ToLower(strings.TrimSpace(cfg.AppLocale))
	switch cfg.AppLocale {
	case "en", "id":
	default:
		return nil, fmt.Errorf("APP_LOCALE must be one of en, id (got %q)", cfg.AppLocale)
	}
	return &cfg, nil
}

func (c *Config) IsDevelopment() bool {
	return c.AppEnv == "development" || c.AppEnv == "local"
}

// loadDotEnv merges .env into the process environment for local runs. Real
// environment variables always win, so deployments can ignore the file.
func loadDotEnv() {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		slog.Warn("could not load .env", "error", err)
	}
}

// DatabaseDSN builds the PostgreSQL DSN from the discrete DB_* settings.
func (c *Config) DatabaseDSN() string {
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.DBUsername, c.DBPassword),
		Host:   net.JoinHostPort(c.DBHost, strconv.Itoa(c.DBPort)),
		Path:   "/" + c.DBDatabase,
	}
	q := u.Query()
	q.Set("sslmode", c.DBSSLMode)
	q.Set("connect_timeout", "3")
	u.RawQuery = q.Encode()
	return u.String()
}
