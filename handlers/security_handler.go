package handlers

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/alexedwards/argon2id"
)

func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func GenerateToken() (string, error) {
	bytes := make([]byte, 32)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func VerifyPassword(password, hash string) (bool, error) {
	checkHash, _, err := argon2id.CheckHash(hash, password)
	if err != nil {
		return false, err
	}
	return checkHash, nil
}
