package lib

import (
	"go-fiber-svelte/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var Jwt = jwtLib{secret: config.APP_SECRET, duration: config.APP_JWT_DURATION}

type jwtLib struct {
	secret   string
	duration string
}

func (j *jwtLib) Create(userId int) (string, error) {
	dur, err := time.ParseDuration(j.duration)
	if err != nil {
		dur = 24 * time.Hour
	}
	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(dur).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secret))
}

func (j *jwtLib) Verify(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		return []byte(j.secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}
	return claims, nil
}
