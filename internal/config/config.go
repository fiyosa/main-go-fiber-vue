package config

import (
	"crypto/sha256"
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
	// AppCORSOrigins is a comma-separated allowlist of browser origins allowed
	// to call the API cross-origin (the Vite dev server runs on another port).
	AppCORSOrigins string `envconfig:"APP_CORS_ORIGINS" default:"http://localhost:5173"`

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

	// Hashid* configure how numeric primary keys are encoded in API payloads.
	// Rotating HASHID_SALT invalidates every id handed out before, so it must be
	// set once and kept. Production must set it explicitly.
	HashidSalt      string `envconfig:"HASHID_SALT"`
	HashidMinLength int    `envconfig:"HASHID_MIN_LENGTH" default:"8"`

	// Storage* configure the S3-compatible object store (Garage in development).
	// The storage layer stays disabled unless endpoint, bucket and both
	// credentials are set, so a missing object store never blocks startup.
	StorageEndpoint  string `envconfig:"STORAGE_ENDPOINT" default:"http://127.0.0.1:3900"`
	StorageRegion    string `envconfig:"STORAGE_REGION" default:"garage"`
	StorageBucket    string `envconfig:"STORAGE_BUCKET" default:"app"`
	StorageAccessKey string `envconfig:"STORAGE_ACCESS_KEY"`
	StorageSecretKey string `envconfig:"STORAGE_SECRET_KEY"`
	StoragePathStyle bool   `envconfig:"STORAGE_PATH_STYLE" default:"true"`

	// MCP exposes a Model Context Protocol server on MCPPath. It is off unless
	// explicitly enabled with MCP_ENABLED=true.
	MCPEnabled        bool          `envconfig:"MCP_ENABLED" default:"false"`
	MCPPath           string        `envconfig:"MCP_PATH" default:"/api/mcp"`
	MCPToken          string        `envconfig:"MCP_TOKEN"`
	MCPScopes         string        `envconfig:"MCP_SCOPES" default:"mcp:read"`
	MCPMaxBodyBytes   int64         `envconfig:"MCP_MAX_BODY_BYTES" default:"1048576"`
	MCPAllowedOrigins string        `envconfig:"MCP_ALLOWED_ORIGINS"`
	MCPSessionTimeout time.Duration `envconfig:"MCP_SESSION_TIMEOUT" default:"60s"`
	MCPStateless      bool          `envconfig:"MCP_STATELESS" default:"true"`
	// MCPSQLMaxRows caps how many rows run_sql returns per call.
	MCPSQLMaxRows int `envconfig:"MCP_SQL_MAX_ROWS" default:"1000"`
	// MCPSQLTimeout bounds how long one run_sql statement may run.
	MCPSQLTimeout time.Duration `envconfig:"MCP_SQL_TIMEOUT" default:"30s"`

	// BasicAuth* configure HTTP Basic authentication used to protect internal
	// endpoints such as the API reference UI.
	BasicAuthUsername string `envconfig:"BASIC_AUTH_USERNAME"`
	BasicAuthPassword string `envconfig:"BASIC_AUTH_PASSWORD"`

	// Cron* configure the cron job engine. Cron is off unless CRON_ENABLED=true.
	// CRON_SECRET_KEY is required when cron is enabled to encrypt credentials.
	CronEnabled  bool   `envconfig:"CRON_ENABLED" default:"false"`
	CronSecretKey string `envconfig:"CRON_SECRET_KEY"`
}

// Load parses the configuration and enforces every runtime secret, so a
// misconfigured application fails fast instead of starting with weak settings.
func Load() (*Config, error) {
	return load(true)
}

// LoadForConsole parses the configuration without enforcing runtime secrets.
// Console commands only introspect the routing table, which does not depend on
// HASHID_SALT, JWT_SECRET, CRON_SECRET_KEY or MCP_TOKEN, so they must stay usable
// in production without those values in the environment.
func LoadForConsole() (*Config, error) {
	return load(false)
}

func load(validateSecrets bool) (*Config, error) {
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
	cfg.MCPPath = normalizePath(cfg.MCPPath)
	if !validateSecrets {
		return &cfg, nil
	}
	if err := cfg.validateMCP(); err != nil {
		return nil, err
	}
	if err := cfg.validateHashid(); err != nil {
		return nil, err
	}
	if err := cfg.validateCron(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// normalizePath trims surrounding spaces and guarantees a leading slash.
func normalizePath(p string) string {
	p = strings.TrimSpace(p)
	if p == "" {
		return "/api/mcp"
	}
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	return p
}

// validateMCP enforces that the static bearer token exists when the endpoint
// is mounted. A missing token is a hard startup error, never a silent bypass.
func (c *Config) validateMCP() error {
	if !c.MCPEnabled {
		return nil
	}
	if strings.TrimSpace(c.MCPToken) == "" {
		return errors.New("MCP_TOKEN is required when MCP_ENABLED=true")
	}
	return nil
}

// validateHashid makes the id salt mandatory in production. Without it every
// hashid would be predictable, so failing fast is safer than warning.
func (c *Config) validateHashid() error {
	if c.IsDevelopment() {
		return nil
	}
	if strings.TrimSpace(c.HashidSalt) == "" {
		return errors.New("HASHID_SALT is required outside development")
	}
	return nil
}

// validateCron enforces that the secret key exists when cron is enabled.
func (c *Config) validateCron() error {
	if !c.CronEnabled {
		return nil
	}
	if strings.TrimSpace(c.CronSecretKey) == "" {
		return errors.New("CRON_SECRET_KEY is required when CRON_ENABLED=true")
	}
	return nil
}

// HashidSaltOrDefault returns the salt used to encode ids. Development falls
// back to a fixed local salt so a fresh checkout runs without extra setup.
func (c *Config) HashidSaltOrDefault() string {
	if salt := strings.TrimSpace(c.HashidSalt); salt != "" {
		return salt
	}
	return "dev-only-hashid-salt"
}

// MCPTokenHash returns the SHA-256 sum of the configured bearer token. Callers
// compare this against the hash of a presented token; the raw token is never
// stored or logged.
func (c *Config) MCPTokenHash() [32]byte {
	return sha256.Sum256([]byte(c.MCPToken))
}

func (c *Config) IsDevelopment() bool {
	return c.AppEnv == "development" || c.AppEnv == "local"
}

// IsStorageConfigured reports whether every required storage setting is
// present. Bucket-only or credential-only setups stay disabled instead of
// failing at request time.
func (c *Config) IsStorageConfigured() bool {
	return strings.TrimSpace(c.StorageEndpoint) != "" &&
		strings.TrimSpace(c.StorageBucket) != "" &&
		strings.TrimSpace(c.StorageAccessKey) != "" &&
		strings.TrimSpace(c.StorageSecretKey) != ""
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
