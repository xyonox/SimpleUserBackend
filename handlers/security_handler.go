package handlers

import (
	"github.com/alexedwards/argon2id"
)

func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func VerifyPassword(password, hash string) (bool, error) {
	checkHash, _, err := argon2id.CheckHash(hash, password)
	if err != nil {
		return false, err
	}
	return checkHash, nil
}
