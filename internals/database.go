package internals

import (
	"errors"
	"fmt"
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

const BucketVersion = "BucketVersion"
const DBVersionName = "DBVersion"

func init() {

}

func CreateDatabaseIfNotExists() {
	var err error
	Database, err = bolt.Open("db_locker.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	err = Database.Update(createBuckets)
	if err != nil {
		log.Fatal(err)
	}
	err = writeVersion()
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
	_, err = tx.CreateBucketIfNotExists([]byte(BucketVersion))
	if err != nil {
		return err
	}
	return nil
}

func writeVersion() error {
	tx, err := Database.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	empty, err := IsBucketEmpty(tx, BucketVersion)
	if empty {
		bucket := tx.Bucket([]byte(BucketVersion))
		if bucket == nil {
			return err
		}
		err = bucket.Put([]byte(DBVersionName), []byte(DBSchemaVersion))
		if err != nil {
			return err
		}
	}
	tx.Commit()
	return nil
}

func getValue(tx *bolt.Tx, key []byte, bucketName string) (value []byte, err error) {
	bucket := tx.Bucket([]byte(bucketName))
	if bucket == nil {
		return nil, errors.New("Bucket: " + bucketName + " doesn't exist")
	}
	value = bucket.Get(key)
	if value == nil {
		return nil, errors.New("No value for id: " + string(key) + " in bucket: " + bucketName + " found")
	}
	return value, nil
}

func saveValue(tx *bolt.Tx, key, value []byte, bucketName string) error {
	bucket := tx.Bucket([]byte(bucketName))
	if bucket == nil {
		return errors.New("Bucket: " + bucketName + " doesn't exist")
	}
	err := bucket.Put(key, value)
	if err != nil {
		return err
	}
	return nil
}
func deleteValue(tx *bolt.Tx, key []byte, bucketName string) error {
	bucket := tx.Bucket([]byte(bucketName))
	if bucket == nil {
		return errors.New("Bucket: " + bucketName + " doesn't exist")
	}
	err := bucket.Delete(key)
	if err != nil {
		return err
	}
	return nil
}
func deleteValues(tx *bolt.Tx, key []byte, bucketNames ...string) error {
	for i := range bucketNames {
		err := deleteValue(tx, key, bucketNames[i])
		if err != nil {
			return err
		}
	}
	return nil
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
func IsDatabaseEmpty() (bool, error) {
	var err error
	if Database == nil {
		fmt.Println("Database is nil")
	}
	tx, err := Database.Begin(false)
	defer tx.Rollback()
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
