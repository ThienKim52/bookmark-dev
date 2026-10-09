package genpass

import (
	"crypto/rand"
	"math/big"
)

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

// GenPass interface for password generation
//
//go:generate mockery --name=GenPass --filename=genpass.go --outpkg=mocks_genpass
type GenPass interface {
	GeneratePassword(passwordLength int) (string, error)
}

type genPass struct{}

// NewGenPass creates a new GenPass instance
func NewGenPass() GenPass {
	return &genPass{}
}

func (g *genPass) GeneratePassword(passwordLength int) (string, error) {
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
