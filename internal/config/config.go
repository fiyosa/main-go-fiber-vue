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

	// MCP exposes a read-only Model Context Protocol server on MCPPath. It is
	// disabled per environment when MCP_ENABLED=false.
	MCPEnabled        bool          `envconfig:"MCP_ENABLED" default:"true"`
	MCPPath           string        `envconfig:"MCP_PATH" default:"/api/mcp"`
	MCPBearerToken    string        `envconfig:"MCP_BEARER_TOKEN"`
	MCPActorUserID    string        `envconfig:"MCP_ACTOR_USER_ID"`
	MCPScopes         string        `envconfig:"MCP_SCOPES" default:"mcp:read"`
	MCPMaxBodyBytes   int64         `envconfig:"MCP_MAX_BODY_BYTES" default:"1048576"`
	MCPAllowedOrigins string        `envconfig:"MCP_ALLOWED_ORIGINS"`
	MCPSessionTimeout time.Duration `envconfig:"MCP_SESSION_TIMEOUT" default:"60s"`
	MCPStateless      bool          `envconfig:"MCP_STATELESS" default:"true"`
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
	cfg.MCPPath = normalizePath(cfg.MCPPath)
	if err := cfg.validateMCP(); err != nil {
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
	if strings.TrimSpace(c.MCPBearerToken) == "" {
		return errors.New("MCP_BEARER_TOKEN is required when MCP_ENABLED=true")
	}
	return nil
}

// MCPTokenHash returns the SHA-256 sum of the configured bearer token. Callers
// compare this against the hash of a presented token; the raw token is never
// stored or logged.
func (c *Config) MCPTokenHash() [32]byte {
	return sha256.Sum256([]byte(c.MCPBearerToken))
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
