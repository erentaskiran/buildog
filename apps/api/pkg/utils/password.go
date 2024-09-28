package utils

import (
	"crypto/rand"
	"math/big"
)

const (
	passwordLength = 12
	charset        = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+[]{}|;:,.<>?/`~"
)

func GeneratePassword(length int) (string, error) {
	password := make([]byte, length)
	for i := range password {
		randomInt, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		password[i] = charset[randomInt.Int64()]
	}
	return string(password), nil
}
