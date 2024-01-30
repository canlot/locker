package internals

import (
	"github.com/boltdb/bolt"
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

const PublicKeyKeyName = "PublicKey"
const PrivateKeyHashKeyName = "PrivateKeyHash"

func init() {
	Database, err := bolt.Open("locker.db", 0600, nil)
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
	return nil
}
