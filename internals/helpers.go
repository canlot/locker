package internals

import (
	"errors"
	bolt "go.etcd.io/bbolt"
	"main/cryptography"
	"time"
)

type DBStore struct {
}

type DataInformation struct {
	Label      string
	CreateTime time.Time
}

type FileInformation struct {
	Path       string
	CreateTime time.Time
}

var FileEncryptionEnding = ".lock"

func getAndDecryptPrivateKey(loginId, passwordHash []byte, tx *bolt.Tx) (privateKey []byte, err error) {
	loginPrivateKeyBucket := tx.Bucket([]byte(BucketLoginPrivateKeyEncrypted))
	privateKeyEncrypted := loginPrivateKeyBucket.Get(loginId)
	if privateKeyEncrypted == nil {
		return nil, errors.New("No private key exists for this login: " + string(loginId))
	}
	decryptedPrivateKey, err := cryptography.DecryptDataSymmetric(passwordHash, privateKeyEncrypted)
	if err != nil {
		return nil, err
	}
	return decryptedPrivateKey, nil
}

func getPrivateKeyHash(tx *bolt.Tx) (hash []byte, err error) {
	privateKeyHashBucket := tx.Bucket([]byte(BucketPrivateKeyHash))
	privateKeyHash := privateKeyHashBucket.Get([]byte(PrivateKeyHashKeyName))
	if privateKeyHash == nil {
		return nil, errors.New("Private key hash is not set")
	}
	return privateKeyHash, nil
}

func getPublicKey(tx *bolt.Tx) (publicKey []byte, err error) {
	publicKeyBucket := tx.Bucket([]byte(BucketPublicKey))
	publicKey = publicKeyBucket.Get([]byte(PublicKeyKeyName))
	if publicKey == nil {
		return nil, errors.New("Public key is empty")
	}
	return publicKey, nil
}

func ThrowErrorIfUidAlreadyExist(tx *bolt.Tx, uid []byte, buckets ...string) error {
	for i := range buckets {
		bucket := tx.Bucket([]byte(buckets[i]))
		if bucket == nil {
			return errors.New("Bucket: " + buckets[i] + " is empty")
		}
		value := bucket.Get(uid)
		if value != nil {
			return errors.New("A glitch in the universe has been happen, uuid already exist, please retry")
		}
	}
	return nil

}
