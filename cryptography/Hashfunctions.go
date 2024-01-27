package cryptography

import (
	"github.com/thanhpk/randstr"
	"golang.org/x/crypto/scrypt"
)

func GenerateUserHash() ([]byte, error) {
	return scrypt.Key([]byte("some password"), nil, 32768, 8, 2, 32)
}

func GenerateRandomString() string {
	return randstr.String(64)
}
