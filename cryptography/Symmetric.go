package cryptography

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"io"
	"os"
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

func EncryptFileSymmetric(password []byte, sourcePath, destinationPath string) error {
	if len(password) != 32 {
		return errors.New("Password not 32 bytes long")
	}
	plainFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer plainFile.Close()
	encryptedFile, err := os.Create(destinationPath)
	if err != nil {
		return err
	}
	defer encryptedFile.Close()
	byteReader := io.Reader(plainFile)
	byteWriter := io.Writer(encryptedFile)
	blockCipher, err := aes.NewCipher(password)
	if err != nil {
		return err
	}
	var iv [aes.BlockSize]byte //no initialization vector needed because every data gets different key
	stream := cipher.NewOFB(blockCipher, iv[:])

	cryptWriter := &cipher.StreamWriter{S: stream, W: byteWriter}

	if _, err = io.Copy(cryptWriter, byteReader); err != nil {
		return err
	}
	return nil
}

func DecryptFileSymmetric(password []byte, sourcePath, destinationPath string) (err error) {
	if len(password) != 32 {
		errors.New("Password not 32 bytes long")
	}
	encryptedFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer encryptedFile.Close()
	decryptedFile, _ := os.Create(destinationPath)
	if err != nil {
		return err
	}
	defer encryptedFile.Close()
	byteReader := io.Reader(encryptedFile)
	byteWriter := io.Writer(decryptedFile)
	blockCipher, err := aes.NewCipher(password)
	if err != nil {
		return err
	}
	var iv [aes.BlockSize]byte
	stream := cipher.NewOFB(blockCipher, iv[:])

	cryptReader := &cipher.StreamReader{S: stream, R: byteReader}
	// Copy the input to the output stream, decrypting as we go.
	if _, err := io.Copy(byteWriter, cryptReader); err != nil {
		return err
	}
	return nil
}
