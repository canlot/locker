package internals

import (
	bolt "go.etcd.io/bbolt"
	"log"
)

var Database *bolt.DB

const BucketLoginInformation = "BucketLoginInformation"
const BucketLoginPrivateKeyEncrypted = "BucketLoginPrivateKeyEncrypted"
const BucketPrivateKeyHash = "BucketPublicKeyHash"
const BucketPublicKey = "BucketPublicKey"
const BucketDataPasswordEncrypted = "BucketDataPasswordEncrypted"
const BucketDataInformation = "BucketDataInformation"
const BucketDataEncrypted = "BucketDataEncrypted"
const BucketDataHash = "BucketDataHash"

const BucketFilePasswordEncrypted = "BucketFilePasswordEncrypted"
const BucketFileInformation = "BucketFileInformation"
const BucketFileHash = "BucketFileHash"

const PublicKeyKeyName = "PublicKey"
const PrivateKeyHashKeyName = "PrivateKeyHash"

func init() {
	var err error
	Database, err = bolt.Open("db_locker.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	err = Database.Update(createBuckets)
	if err != nil {
		log.Fatal(err)
	}
}

func createBuckets(tx *bolt.Tx) error {
	var err error
	_, err = tx.CreateBucketIfNotExists([]byte(BucketLoginInformation))
	if err != nil {
		return err
	}
	_, err = tx.CreateBucketIfNotExists([]byte(BucketLoginPrivateKeyEncrypted))
	if err != nil {
		return err
	}
	_, err = tx.CreateBucketIfNotExists([]byte(BucketPrivateKeyHash))
	if err != nil {
		return err
	}
	_, err = tx.CreateBucketIfNotExists([]byte(BucketPublicKey))
	if err != nil {
		return err
	}
	_, err = tx.CreateBucketIfNotExists([]byte(BucketDataPasswordEncrypted))
	if err != nil {
		return err
	}
	_, err = tx.CreateBucketIfNotExists([]byte(BucketDataInformation))
	if err != nil {
		return err
	}
	_, err = tx.CreateBucketIfNotExists([]byte(BucketDataEncrypted))
	if err != nil {
		return err
	}
	_, err = tx.CreateBucketIfNotExists([]byte(BucketDataHash))
	if err != nil {
		return err
	}
	_, err = tx.CreateBucketIfNotExists([]byte(BucketFilePasswordEncrypted))
	if err != nil {
		return err
	}
	_, err = tx.CreateBucketIfNotExists([]byte(BucketFileInformation))
	if err != nil {
		return err
	}
	_, err = tx.CreateBucketIfNotExists([]byte(BucketFileHash))
	if err != nil {
		return err
	}
	return nil
}
