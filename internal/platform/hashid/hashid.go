// Package hashid converts numeric primary keys into opaque strings at the HTTP
// boundary. It is obfuscation, not encryption: authorization is still enforced
// by the permission middleware.
package hashid

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/speps/go-hashids/v2"
)

var (
	// ErrInvalidID is returned when a string is not a single valid hashid.
	ErrInvalidID = errors.New("hashid: invalid id")
	// ErrUnavailable is returned when the encoder was not configured.
	ErrUnavailable = errors.New("hashid: encoder unavailable")
)

// Encoder encodes and decodes single numeric ids.
type Encoder struct {
	h *hashids.HashID
}

// New builds an encoder. The salt must be non-empty; rotating it invalidates
// every id that was handed out before.
func New(salt string, minLength int) (*Encoder, error) {
	if strings.TrimSpace(salt) == "" {
		return nil, errors.New("hashid: salt is required")
	}
	if minLength < 4 {
		minLength = 4
	}

	data := hashids.NewData()
	data.Salt = salt
	data.MinLength = minLength

	h, err := hashids.NewWithData(data)
	if err != nil {
		return nil, fmt.Errorf("hashid: init encoder: %w", err)
	}
	return &Encoder{h: h}, nil
}

// Encode turns an id into its hashid form.
func (e *Encoder) Encode(id uint64) (string, error) {
	if e == nil || e.h == nil {
		return "", ErrUnavailable
	}
	if id > math.MaxInt64 {
		return "", ErrInvalidID
	}
	return e.h.EncodeInt64([]int64{int64(id)})
}

// Decode parses a hashid back into an id.
func (e *Encoder) Decode(raw string) (uint64, error) {
	if e == nil || e.h == nil {
		return 0, ErrUnavailable
	}
	if strings.TrimSpace(raw) == "" {
		return 0, ErrInvalidID
	}

	ids, err := e.h.DecodeInt64WithError(raw)
	if err != nil || len(ids) != 1 || ids[0] < 0 {
		return 0, ErrInvalidID
	}
	return uint64(ids[0]), nil
}
