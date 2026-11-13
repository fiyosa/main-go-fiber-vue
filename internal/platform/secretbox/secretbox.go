package secretbox

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

var (
	ErrInvalidCiphertext = errors.New("invalid ciphertext format")
	ErrEmptyKey          = errors.New("secret key cannot be empty")
)

// deriveKey derives a 32-byte AES key from the secret string using SHA-256.
func deriveKey(secretKey string) ([]byte, error) {
	if secretKey == "" {
		return nil, ErrEmptyKey
	}
	hash := sha256.Sum256([]byte(secretKey))
	return hash[:], nil
}

// Encrypt encrypts plaintext using AES-256-GCM with secretKey.
// Returns a standard base64 string containing [nonce + ciphertext + tag].
func Encrypt(plaintext, secretKey string) (string, error) {
	key, err := deriveKey(secretKey)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("create gcm: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts standard base64 encoded [nonce + ciphertext + tag] using secretKey.
func Decrypt(encodedCiphertext, secretKey string) (string, error) {
	key, err := deriveKey(secretKey)
	if err != nil {
		return "", err
	}

	data, err := base64.StdEncoding.DecodeString(encodedCiphertext)
	if err != nil {
		return "", fmt.Errorf("%w: base64 decode failed", ErrInvalidCiphertext)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("create gcm: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", ErrInvalidCiphertext
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("%w: decryption failed", ErrInvalidCiphertext)
	}

	return string(plaintext), nil
}
