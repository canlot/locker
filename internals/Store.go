package internals

import (
	"errors"
	"github.com/boltdb/bolt"
	"golang.org/x/crypto/scrypt"
	"time"
)

type DBStore struct {
}

var Store DBStore

func (s DBStore) IsDatabaseEmpty() (bool, error) {
	var err error
	tx, err := Database.Begin(false)
	tx.Commit()
	if err != nil {
		return false, err
	}
	var empty bool
	empty, err = IsBucketEmpty(tx, BucketLoginInformation)
	if err != nil {
		return false, err
	}
	if !empty {
		return false, nil
	}
	empty, err = IsBucketEmpty(tx, BucketLoginPrivateKeyEncrypted)
	if err != nil {
		return false, err
	}
	if !empty {
		return false, nil
	}
	empty, err = IsBucketEmpty(tx, BucketLoginPrivateKeyEncrypted)
	if err != nil {
		return false, err
	}
	if !empty {
		return false, nil
	}

	return true, nil
}
func IsBucketEmpty(tx *bolt.Tx, bucketName string) (bool, error) {
	bucket := tx.Bucket([]byte(bucketName))
	if bucket == nil {
		return false, errors.New("Bucket not created: " + bucketName)
	}
	c := bucket.Cursor()
	first, _ := c.First()
	if first != nil {
		return false, nil
	}
	return true, nil
}
func CreateFirstLoginWithRSAKeys(username, password string) error {
	createTime := time.Now()
	passwordHash, err := scrypt.Key([]byte(password), nil, 32768, 8, 1, 32)
	if err != nil {
		return err
	}
	u := uuid.New()
}
