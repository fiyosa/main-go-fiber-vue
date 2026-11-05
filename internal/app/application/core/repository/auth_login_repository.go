package repository

import (
	"context"
	"encoding/json"
	"time"

	"gorm.io/gorm"

	"go-fiber-vue/internal/app/domain/core/entity"
	"go-fiber-vue/internal/app/infrastructure/support/exceptions"
	"go-fiber-vue/internal/platform/auth"
)

// TokenService defines the JWT operations required by the login repository.
// It is injected by the bootstrap so the repository stays decoupled from auth/platform.
type TokenService interface {
	IssueAccessToken(userID uint64, username string, tokenVersion int) (string, time.Time, error)
	NewRefreshToken() (raw string, hash string, err error)
	AccessTTL() time.Duration
	RefreshTTL() time.Duration
}

// Meta carries client metadata recorded on each login.
type Meta struct {
	IP        string
	UserAgent string
}

// LoginResult is returned on a successful login.
type LoginResult struct {
	Token        string
	RefreshToken string
	ExpiresIn    int64
}

// AuthLoginRepository is the single interface for the POST /api/v1/auth/login
// endpoint. It owns all business logic: credential check, JWT issue, and session
// persistence. One file per API — no use-case layer.
type AuthLoginRepository interface {
	Login(ctx context.Context, identifier, password string, meta Meta) (*LoginResult, error)
}

// gormAuthLoginRepository is the GORM implementation of AuthLoginRepository.
type gormAuthLoginRepository struct {
	db     *gorm.DB
	users  UserRepository
	auths  AuthRepository
	tokens TokenService
}

// NewAuthLoginRepository wires the login repository with its GORM implementation.
// It requires the generic user and auth repositories so it can compose them.
func NewAuthLoginRepository(db *gorm.DB, users UserRepository, auths AuthRepository, tokens TokenService) AuthLoginRepository {
	return &gormAuthLoginRepository{db: db, users: users, auths: auths, tokens: tokens}
}

// Login implements AuthLoginRepository.
func (r *gormAuthLoginRepository) Login(ctx context.Context, identifier, password string, meta Meta) (*LoginResult, error) {
	if r.users == nil || r.auths == nil || r.tokens == nil {
		return nil, exceptions.ServiceUnavailable(nil)
	}

	user, err := r.users.FindByUsernameOrEmail(ctx, identifier)
	if err != nil {
		return nil, err
	}
	if user == nil || !auth.VerifyPassword(user.Password, password) {
		return nil, exceptions.Unauthorized("lang.invalid_credentials", nil)
	}

	accessToken, accessExpiresAt, err := r.tokens.IssueAccessToken(user.ID, user.Username, user.TokenVersion)
	if err != nil {
		return nil, exceptions.Internal(err)
	}

	refreshRaw, refreshHash, err := r.tokens.NewRefreshToken()
	if err != nil {
		return nil, exceptions.Internal(err)
	}

	now := time.Now().UTC()
	row := &entity.Auth{
		UserID:           user.ID,
		Token:            accessToken,
		RefreshToken:     refreshHash,
		TokenExpiresAt:   accessExpiresAt,
		RefreshExpiresAt: now.Add(r.tokens.RefreshTTL()),
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if meta.IP != "" || meta.UserAgent != "" {
		raw, marshalErr := json.Marshal(map[string]string{
			"ip":         meta.IP,
			"user_agent": meta.UserAgent,
		})
		if marshalErr == nil {
			value := string(raw)
			row.Meta = &value
		}
	}

	if err := r.auths.CreateSession(ctx, row, now); err != nil {
		return nil, err
	}

	return &LoginResult{
		Token:        accessToken,
		RefreshToken: refreshRaw,
		ExpiresIn:    int64(r.tokens.AccessTTL().Seconds()),
	}, nil
}
