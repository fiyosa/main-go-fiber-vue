package auth

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v3"

	"go-fiber-vue/internal/app/application/core/repository/user"
	"go-fiber-vue/internal/app/domain/core/entity"
	"go-fiber-vue/internal/app/infrastructure/support/exceptions"
	platformauth "go-fiber-vue/internal/platform/auth"
)

// Auth cookie names.
const (
	AccessCookieName  = "access_token"
	RefreshCookieName = "refresh_token"
)

// SetAuthCookies stores the issued access and refresh tokens in HTTP-only
// cookies. Cookies are not readable by JavaScript and never appear in the
// response body.
func SetAuthCookies(fctx fiber.Ctx, accessToken, refreshToken string, accessTTL, refreshTTL time.Duration, secure bool) {
	fctx.Cookie(&fiber.Cookie{
		Name:     AccessCookieName,
		Value:    accessToken,
		MaxAge:   int(accessTTL.Seconds()),
		HTTPOnly: true,
		Secure:   secure,
		SameSite: fiber.CookieSameSiteLaxMode,
	})
	fctx.Cookie(&fiber.Cookie{
		Name:     RefreshCookieName,
		Value:    refreshToken,
		MaxAge:   int(refreshTTL.Seconds()),
		HTTPOnly: true,
		Secure:   secure,
		SameSite: fiber.CookieSameSiteLaxMode,
	})
}

// ClearAuthCookies clears both auth cookies by setting MaxAge to -1.
func ClearAuthCookies(fctx fiber.Ctx, secure bool) {
	fctx.Cookie(&fiber.Cookie{
		Name:     AccessCookieName,
		Value:    "",
		MaxAge:   -1,
		HTTPOnly: true,
		Secure:   secure,
		SameSite: fiber.CookieSameSiteLaxMode,
	})
	fctx.Cookie(&fiber.Cookie{
		Name:     RefreshCookieName,
		Value:    "",
		MaxAge:   -1,
		HTTPOnly: true,
		Secure:   secure,
		SameSite: fiber.CookieSameSiteLaxMode,
	})
}

// TokenService defines the JWT operations required by the login repository.
// It is injected by the bootstrap so the repository stays decoupled from auth/platform.
type TokenService interface {
	IssueAccessToken(userID uint64, email string) (string, error)
	NewRefreshToken() (raw string, hash string, err error)
	AccessTTL() time.Duration
	RefreshTTL() time.Duration
}

// Meta carries client metadata recorded on each login.
type Meta struct {
	IP        string
	UserAgent string
}

// LoginResult is returned on a successful login. The TTLs are exposed so the
// presentation layer can size the auth cookies without importing the token
// service.
type LoginResult struct {
	Token            string
	RefreshToken     string
	ExpiresIn        int64
	RefreshExpiresIn int64
}

// LoginRepository is the single interface for the POST /api/v1/auth/login
// endpoint. It owns all business logic: credential check, JWT issue, and session
// persistence. One file per API — no use-case layer.
type LoginRepository interface {
	Login(ctx context.Context, identifier, password string, meta Meta) (*LoginResult, error)
}

// gormLoginRepository is the GORM implementation of LoginRepository.
type gormLoginRepository struct {
	users  user.Repository
	auths  Repository
	tokens TokenService
}

// NewLoginRepository wires the login repository with its GORM implementation.
// It composes the generic user and auth repositories.
func NewLoginRepository(users user.Repository, auths Repository, tokens TokenService) LoginRepository {
	return &gormLoginRepository{users: users, auths: auths, tokens: tokens}
}

// Login implements LoginRepository.
func (r *gormLoginRepository) Login(ctx context.Context, email, password string, meta Meta) (*LoginResult, error) {
	if r.users == nil || r.auths == nil || r.tokens == nil {
		return nil, exceptions.ServiceUnavailable(nil)
	}

	account, err := r.users.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if account == nil || !platformauth.VerifyPassword(account.Password, password) {
		return nil, exceptions.Unauthorized("lang.invalid_credentials", nil)
	}

	accessToken, err := r.tokens.IssueAccessToken(account.ID, account.Email)
	if err != nil {
		return nil, exceptions.Internal(err)
	}

	refreshRaw, refreshHash, err := r.tokens.NewRefreshToken()
	if err != nil {
		return nil, exceptions.Internal(err)
	}

	now := time.Now().UTC()
	row := &entity.Auth{
		UserID:       account.ID,
		Token:        accessToken,
		RefreshToken: refreshHash,
		CreatedAt:    now,
		UpdatedAt:    now,
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

	if err := r.auths.CreateSession(ctx, row); err != nil {
		return nil, err
	}

	return &LoginResult{
		Token:            accessToken,
		RefreshToken:     refreshRaw,
		ExpiresIn:        int64(r.tokens.AccessTTL().Seconds()),
		RefreshExpiresIn: int64(r.tokens.RefreshTTL().Seconds()),
	}, nil
}
