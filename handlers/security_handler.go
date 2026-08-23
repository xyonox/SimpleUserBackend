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
	index := 0
	maxSize := len(salt)

	hashedPassword := make([]byte, maxSize)

	for _, saltedPasswordChar := range saltedPassword {
		var finalHash byte
		finalHash = saltedPasswordChar

		for j := range states {
			stateChar := states[j]
			if j == 0 {
				finalHash += states[j+1] + states[len(states)-1]
				states[0] = states[j+1] + states[len(states)-1]*states[j+5]
			} else if j == len(states)-1 {
				finalHash += states[0] + states[j-1]*states[j-5]
			} else {
				finalHash += states[j+1] + states[j-1]
			}
			if stateChar == finalHash {
				val := finalHash*stateChar + 4
				if val > 127 {
					finalHash = val - 120
					states[j] = val ^ stateChar
				} else if val%2 == 0 {
					finalHash = val - stateChar
					states[j] = val + stateChar
				} else if (val-stateChar)%2 == 0 {
					finalHash = val + stateChar*4
					states[j] = val - (val ^ (stateChar * 22))
				} else {
					finalHash = val ^ stateChar
					states[j] = val + (val ^ (stateChar * 2))
				}
			} else if finalHash > stateChar {
				val := finalHash - stateChar
				if val > 127 {
					finalHash = val ^ stateChar
					states[j] = val - (stateChar * 30)
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
					states[j] = (val - stateChar) * 5
				} else {
					finalHash = val ^ stateChar
					states[j] = (val + stateChar) * (val ^ (stateChar * 20))
				}
			}
		}
		if index == maxSize {
			index = 0
		}
		hashedPassword[index] += finalHash
		index++
	}

	return hashedPassword, nil
}

func HashPassword(password string) ([]byte, []byte, error) {
	salt := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, nil, err
	}
	hashedPassword, err := HashPasswordWithSalt(password, salt)

	return hashedPassword, salt, err
}

func VerifyPasswordHash(password string, hashed []byte, salt []byte) bool {
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
