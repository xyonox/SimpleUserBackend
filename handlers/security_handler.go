package handlers

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
)

func HashPasswordWithSalt(password string, salt []byte) ([]byte, error) {

	saltedPassword := append(salt, []byte(password)...)

	var states []byte
	for _, saltChar := range saltedPassword {
		states = append(states, saltChar)
	}

	var hashedPassword []byte

	for _, saltedPasswordChar := range saltedPassword {
		var finalHash byte
		finalHash = saltedPasswordChar
		for j := range states {
			stateChar := states[j]
			if stateChar == finalHash {
				val := finalHash*stateChar + 4
				if val > 127 {
					finalHash = val - 120
					states[j] = val ^ stateChar
				} else if val%2 == 0 {
					finalHash = val - stateChar
					states[j] = val + stateChar
				} else {
					finalHash = val ^ stateChar
					states[j] = val + (val ^ (stateChar * 2))
				}
			} else if finalHash > stateChar {
				val := finalHash - stateChar
				if val > 127 {
					finalHash = val ^ stateChar
					states[j] = val - stateChar
				} else if val%2 == 0 {
					finalHash = val + stateChar
					states[j] = (val + stateChar) * 2
				} else {
					finalHash = val ^ stateChar
					states[j] = (val - stateChar) * (val ^ (stateChar * 2))
				}
			} else if finalHash < stateChar {
				val := stateChar - finalHash
				if val > 127 {
					finalHash = val ^ stateChar
					states[j] = val + stateChar
				} else if val%2 == 0 {
					finalHash = val - stateChar
					states[j] = (val - stateChar) * 2
				} else {
					finalHash = val ^ stateChar
					states[j] = (val + stateChar) * (val ^ (stateChar * 2))
				}
			}
		}
		hashedPassword = append(hashedPassword, finalHash)
	}

	return hashedPassword, nil
}

func HashPassword(password string) ([]byte, []byte, error) {
	salt := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, nil, err
	}

	hashedPassword, err := HashPasswordWithSalt(password, salt)

	return hashedPassword, salt, err
}

func CheckPasswordHash(password string, hashed []byte, salt []byte) bool {
	hashedCheck, err := HashPasswordWithSalt(password, salt)
	if err != nil {
		fmt.Println("Error hashing password: ", err)
		return false
	}
	if bytes.Equal(hashed, hashedCheck) {
		return true
	}

	return false
}
