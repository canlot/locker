package cryptography

import (
	"crypto/sha256"
	"github.com/thanhpk/randstr"
	"golang.org/x/crypto/scrypt"
	"io"
	"os"
)

func GenerateUserHash(password []byte) ([]byte, error) {
	return scrypt.Key(password, nil, 32768, 8, 2, 32)
}

func GenerateRandomBytes() []byte {
	return randstr.Bytes(32)
}
func GetSha256Hash(input []byte) (hash []byte, err error) {
	h := sha256.New()
	_, err = h.Write(input)
	if err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}
func GetSha256HashFile(filePath string) (hash []byte, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}
