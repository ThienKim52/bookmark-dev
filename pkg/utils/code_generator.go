package utils

import (
	"crypto/rand"
	"math/big"
)

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

type KeyGenerator interface {
	KeyGenerator(passwordLength int) (string, error)
}

type keyGenerator struct{}

// NewGenPass creates a new GenPass instance
func NewKeyGenerator() KeyGenerator {
	return &keyGenerator{}
}

func (g *keyGenerator) KeyGenerator(passwordLength int) (string, error) {
	password := make([]byte, passwordLength)

	for i := range password {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		password[i] = charset[index.Int64()]
	}

	return string(password), nil
}
