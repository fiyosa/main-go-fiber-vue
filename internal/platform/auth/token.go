package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken  = errors.New("invalid token")
	ErrMissingSecret = errors.New("JWT_SECRET is required")
)

// Fixed claim values. They are not configurable through the environment.
const (
	DefaultIssuer   = "app-api"
	DefaultAudience = "app-spa"

	recommendedSecretLength = 32
)

type Config struct {
	Secret     string
	Issuer     string
	Audience   string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

type Service struct {
	issuer     string
	audience   string
	accessTTL  time.Duration
	refreshTTL time.Duration
	secret     []byte
}

type Claims struct {
	Username string `json:"username"`
	Version  int    `json:"ver"`
	jwt.RegisteredClaims
}

func NewService(cfg Config) (*Service, error) {
	if cfg.Secret == "" {
		return nil, ErrMissingSecret
	}
	if len(cfg.Secret) < recommendedSecretLength {
		slog.Warn("JWT_SECRET is shorter than the recommended 32 characters")
	}

	issuer := cfg.Issuer
	if issuer == "" {
		issuer = DefaultIssuer
	}
	audience := cfg.Audience
	if audience == "" {
		audience = DefaultAudience
	}

	return &Service{
		issuer:     issuer,
		audience:   audience,
		accessTTL:  cfg.AccessTTL,
		refreshTTL: cfg.RefreshTTL,
		secret:     []byte(cfg.Secret),
	}, nil
}

func (s *Service) AccessTTL() time.Duration {
	return s.accessTTL
}

func (s *Service) RefreshTTL() time.Duration {
	return s.refreshTTL
}

func (s *Service) NewRefreshToken() (string, string, error) {
	return NewRefreshToken()
}

func (s *Service) IssueAccessToken(userID uint64, username string, tokenVersion int) (string, time.Time, error) {
	now := time.Now().UTC()
	expiresAt := now.Add(s.accessTTL)
	claims := Claims{
		Username: username,
		Version:  tokenVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Audience:  jwt.ClaimStrings{s.audience},
			Subject:   strconv.FormatUint(userID, 10),
			ID:        newJTI(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.secret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("sign access token: %w", err)
	}
	return signed, expiresAt, nil
}

func (s *Service) ValidateAccessToken(raw string, currentVersion int) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(raw, claims, func(t *jwt.Token) (any, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, ErrInvalidToken
		}
		return s.secret, nil
	}, jwt.WithIssuer(s.issuer), jwt.WithAudience(s.audience), jwt.WithExpirationRequired())
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}
	if claims.Version != currentVersion {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

func NewRefreshToken() (raw string, hash string, err error) {
	buf := make([]byte, 32)
	if _, err = rand.Read(buf); err != nil {
		return "", "", fmt.Errorf("generate refresh token: %w", err)
	}
	raw = base64.RawURLEncoding.EncodeToString(buf)
	return raw, HashRefreshToken(raw), nil
}

func HashRefreshToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func newJTI() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return strconv.FormatInt(time.Now().UnixNano(), 10)
	}
	return hex.EncodeToString(buf)
}
