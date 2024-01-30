package cryptography

import (
	"crypto/sha256"
	"github.com/thanhpk/randstr"
	"golang.org/x/crypto/scrypt"
)

func GenerateUserHash(password []byte) ([]byte, error) {
	return scrypt.Key(password, nil, 32768, 8, 2, 32)
}

func GenerateRandomString() string {
	return randstr.String(64)
}

func GetSha256Hash(input []byte) (hash []byte, err error) {
	h := sha256.New()
	_, err = h.Write(input)
	if err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}
