package cryptography

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"errors"
	"io"
)

func EncryptDataSymmetric(password, plainData []byte) (encryptedData []byte, err error) {
	if len(password) != 32 {
		return nil, errors.New("Password not 32 bytes long")
	}
	byteReader := bytes.NewReader(plainData)
	blockCipher, err := aes.NewCipher(password)
	if err != nil {
		return nil, err
	}
	var iv [aes.BlockSize]byte //no initialization vector needed because every data gets different key
	stream := cipher.NewOFB(blockCipher, iv[:])

	var out bytes.Buffer
	cryptWriter := &cipher.StreamWriter{S: stream, W: &out}
	if _, err := io.Copy(cryptWriter, byteReader); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func DecryptDataSymmetric(password, encryptedData []byte) (plainData []byte, err error) {
	if len(password) != 32 {
		return nil, errors.New("Password not 32 bytes long")
	}
	byteReader := bytes.NewReader(encryptedData)
	blockCipher, err := aes.NewCipher(password)
	if err != nil {
		return nil, err
	}
	var iv [aes.BlockSize]byte
	stream := cipher.NewOFB(blockCipher, iv[:])

	var out bytes.Buffer
	byteWriter := io.Writer(&out)
	cryptReader := &cipher.StreamReader{S: stream, R: byteReader}
	// Copy the input to the output stream, decrypting as we go.
	if _, err := io.Copy(byteWriter, cryptReader); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func EncryptFileSymmetricWithHash(password []byte, fileReader io.Reader, fileWriter, progressWriter io.Writer) ([]byte, error) {
	if len(password) != 32 {
		return nil, errors.New("Password not 32 bytes long")
	}

	blockCipher, err := aes.NewCipher(password)
	if err != nil {
		return nil, err
	}
	var iv [aes.BlockSize]byte //no initialization vector needed because every data gets different key
	stream := cipher.NewOFB(blockCipher, iv[:])

	cryptWriter := &cipher.StreamWriter{S: stream, W: fileWriter}
	hash := sha256.New()
	multiWriter := io.MultiWriter(cryptWriter, hash, progressWriter)

	if _, err = io.Copy(multiWriter, fileReader); err != nil {
		return nil, err
	}
	return hash.Sum(nil), nil
}

func DecryptFileSymmetricWithHash(password []byte, fileReader io.Reader, fileWriter, progressWriter io.Writer) ([]byte, error) {
	if len(password) != 32 {
		errors.New("Password not 32 bytes long")
	}
	blockCipher, err := aes.NewCipher(password)
	if err != nil {
		return nil, err
	}
	var iv [aes.BlockSize]byte
	stream := cipher.NewOFB(blockCipher, iv[:])

	cryptReader := &cipher.StreamReader{S: stream, R: fileReader}

	hash := sha256.New()

	multiWriter := io.MultiWriter(fileWriter, hash, progressWriter)

	// Copy the input to the output stream, decrypting as we go.
	if _, err := io.Copy(multiWriter, cryptReader); err != nil {
		return nil, err
	}
	return hash.Sum(nil), nil
}
