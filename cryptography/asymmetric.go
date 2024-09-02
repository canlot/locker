package cryptography

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
)

func GenerateRSAKeys() (privateKeyBytes, publicKeyBytes []byte, err error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil, err
	}
	err = privateKey.Validate()
	if err != nil {
		return nil, nil, err
	}
	publicKey := privateKey.PublicKey
	privateKeyBytes = x509.MarshalPKCS1PrivateKey(privateKey)
	publicKeyBytes = x509.MarshalPKCS1PublicKey(&publicKey)
	return privateKeyBytes, publicKeyBytes, nil
}

func EncryptDataAsymmetric(publicKeyBytes, dataIn []byte) (encryptedDataOut []byte, err error) {
	publicKey, err := x509.ParsePKCS1PublicKey(publicKeyBytes)
	if err != nil {
		return nil, err
	}
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, dataIn, nil)
	if err != nil {
		return nil, err
	}
	return ciphertext, nil
}
func DecryptDataAsymmetric(privateKeyBytes, encryptedDataIn []byte) (dataOut []byte, err error) {
	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBytes)
	if err != nil {
		return nil, err
	}
	plaintext, err := rsa.DecryptOAEP(sha256.New(), nil, privateKey, encryptedDataIn, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}
