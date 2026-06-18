package lib

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"go-fiber-svelte/internal/config"
	"io"

	"github.com/speps/go-hashids/v2"
	"golang.org/x/crypto/bcrypt"
)

var Hash hash

type hash struct{}

func (*hash) Create(data string) (string, error) {
	result, err := bcrypt.GenerateFromPassword([]byte(data), 12)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func (*hash) Verify(check string, original string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(original), []byte(check)); err != nil {
		return false
	}
	return true
}

func (*hash) EncodeId(data int) (string, error) {
	h := setupHD()
	encode, err := h.Encode([]int{data})
	if err != nil {
		return "", err
	}
	return encode, err
}

func (*hash) DecodeId(data string) (int, error) {
	h := setupHD()
	decode, err := h.DecodeWithError(data)
	if err != nil {
		return -1, err
	}
	return decode[0], err
}

func EncodeStr(data string) (string, error) {
	key := deriveKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, []byte(data), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func DecodeStr(encode string) (string, error) {
	key := deriveKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	ciphertext, err := base64.StdEncoding.DecodeString(encode)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func deriveKey() []byte {
	h := sha256.Sum256([]byte(config.APP_SECRET))
	return h[:]
}

func setupHD() *hashids.HashID {
	hd := hashids.NewData()
	hd.Salt = config.APP_SECRET
	hd.MinLength = 10
	h, _ := hashids.NewWithData(hd)
	return h
}
