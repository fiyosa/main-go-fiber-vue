// Package storage talks to an S3-compatible object store. It knows nothing
// about HTTP, authorization, or the database: callers hand it a path and a
// body, and it never exposes a public URL.
package storage

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"errors"
	"io"
	"strings"
	"unicode"
)

// ErrObjectNotFound reports that no object exists at the requested path.
var ErrObjectNotFound = errors.New("storage: object not found")

// DefaultContentType is used when the caller does not supply a MIME type.
const DefaultContentType = "application/octet-stream"

// keyBytes is the amount of entropy behind every object key (80 bits).
const keyBytes = 10

// keyEncoding produces 16 characters for 10 random bytes. Padding is removed
// and the result is lowercased so keys stay URL-safe without encoding.
var keyEncoding = base32.StdEncoding.WithPadding(base32.NoPadding)

// ObjectStore is the minimal contract the repository needs from a bucket.
type ObjectStore interface {
	Save(ctx context.Context, in SaveInput) (objectPath string, size int64, err error)
	Read(ctx context.Context, objectPath string) (io.ReadCloser, error)
	Delete(ctx context.Context, objectPath string) error
}

// SaveInput describes one object write.
type SaveInput struct {
	Path        string
	ContentType string
	Body        io.Reader
}

// NewKey returns a random lowercase base32 key of 16 characters derived from
// 10 bytes of crypto/rand.
func NewKey() (string, error) {
	buf := make([]byte, keyBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return strings.ToLower(keyEncoding.EncodeToString(buf)), nil
}

// NormalizeFolder validates a caller-supplied folder and returns it without
// surrounding slashes. It rejects empty folders, "."/"..", backslashes, empty
// path segments, and control characters so a value can never escape or confuse
// the intended bucket prefix.
func NormalizeFolder(folder string) (string, error) {
	folder = strings.TrimSpace(folder)
	if folder == "" {
		return "", errors.New("folder is required")
	}
	if strings.ContainsRune(folder, '\\') {
		return "", errors.New("folder must not contain backslashes")
	}
	for _, r := range folder {
		if unicode.IsControl(r) {
			return "", errors.New("folder must not contain control characters")
		}
	}

	folder = strings.Trim(folder, "/")
	if folder == "" {
		return "", errors.New("folder is required")
	}
	for _, segment := range strings.Split(folder, "/") {
		switch segment {
		case "":
			return "", errors.New("folder must not contain empty segments")
		case ".", "..":
			return "", errors.New("folder must not contain relative segments")
		}
	}
	return folder, nil
}
