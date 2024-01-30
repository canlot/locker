package internals

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	bolt "go.etcd.io/bbolt"
	"main/cryptography"
	"time"
)

type DBStore struct {
}

var Store DBStore

func IsDatabaseEmpty() (bool, error) {
	fmt.Println("IsDatabaseEmpty")
	var err error
	if Database == nil {
		fmt.Println("Database is nil")
	}
	tx, err := Database.Begin(false)
	defer fmt.Println("IsDatabaseEmpty Commited")
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
	fmt.Println("CreateFirstLoginWithRSAKeys")
	createTime := time.Now()
	passwordHash, err := cryptography.GenerateUserHash([]byte(password))
	if err != nil {
		return err
	}
	u := uuid.New()
	uuidStringBytes, err := u.MarshalText()
	if err != nil {
		return err
	}
	login := Login{Login: username, CreateTime: createTime}
	loginBytes, err := json.Marshal(login)
	if err != nil {
		return err
	}
	privateKey, publicKey, err := cryptography.GenerateRSAKeys()
	fmt.Printf("private key: %x\n", privateKey)
	fmt.Printf("public key: %x\n", publicKey)
	if err != nil {
		return err
	}
	privateKeyHash, err := cryptography.GetSha256Hash(privateKey)
	fmt.Printf("private key hash: %x\n", privateKeyHash)
	if err != nil {
		return err
	}
	privateKeyEncrypted, err := cryptography.EncryptDataSymmetric(passwordHash, privateKey)
	if err != nil {
		return err
	}
	tx, err := Database.Begin(true)
	if err != nil {
		return err
	}
	userInfoBucket := tx.Bucket([]byte(BucketLoginInformation))
	err = userInfoBucket.Put(uuidStringBytes, loginBytes)
	if err != nil {
		tx.Rollback()
		return err
	}
	userPrivateKeyEncryptedBucket := tx.Bucket([]byte(BucketLoginPrivateKeyEncrypted))
	err = userPrivateKeyEncryptedBucket.Put(uuidStringBytes, privateKeyEncrypted)
	if err != nil {
		tx.Rollback()
		return err
	}
	privateKeyHashBucket := tx.Bucket([]byte(BucketPrivateKeyHash))
	err = privateKeyHashBucket.Put([]byte(PrivateKeyHashKeyName), privateKeyHash)
	if err != nil {
		tx.Rollback()
		return err
	}
	publicKeyBucket := tx.Bucket([]byte(BucketPublicKey))
	err = publicKeyBucket.Put([]byte(PublicKeyKeyName), publicKey)
	if err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}
func CreateLoginWithExistingRSAKeys(existingLogin, existingLoginPassword, newLogin, newLoginPassword string) error {
	fmt.Println("CreateLoginWithExistingRSAKeys")
	tx, err := Database.Begin(true)
	if err != nil {
		tx.Rollback()
		return err
	}
	passwordHashExistingLogin, err := cryptography.GenerateUserHash([]byte(existingLoginPassword))
	if err != nil {
		tx.Rollback()
		return err
	}
	passwordHashNewLogin, err := cryptography.GenerateUserHash([]byte(newLoginPassword))
	if err != nil {
		tx.Rollback()
		return err
	}
	loginInformationBucket := tx.Bucket([]byte(BucketLoginInformation))
	c := loginInformationBucket.Cursor()

	var userid []byte
	for k, v := c.First(); k != nil; k, v = c.Next() {
		var login Login
		err = json.Unmarshal(v, &login)
		if err != nil {
			tx.Rollback()
			return err
		}
		if login.Login == existingLogin {
			userid = k
		}
	}
	if userid == nil {
		tx.Rollback()
		return errors.New("No existing login found: " + existingLogin)
	}
	loginPrivateKeyBucket := tx.Bucket([]byte(BucketLoginPrivateKeyEncrypted))
	privateKeyEncrypted := loginPrivateKeyBucket.Get(userid)
	if privateKeyEncrypted == nil {
		tx.Rollback()
		return errors.New("No private key exists for this login: " + existingLogin)
	}
	privateKeyHashBucket := tx.Bucket([]byte(BucketPrivateKeyHash))
	privateKeyHash := privateKeyHashBucket.Get([]byte(PrivateKeyHashKeyName))
	if privateKeyHash == nil {
		tx.Rollback()
		return errors.New("Private key hash is not set")
	}

	decryptedPrivateKey, err := cryptography.DecryptDataSymmetric([]byte(passwordHashExistingLogin), privateKeyEncrypted)
	fmt.Printf("private key: %x\n", decryptedPrivateKey)

	if err != nil {
		tx.Rollback()
		return err
	}
	generatedPrivateKeyHash, err := cryptography.GetSha256Hash(decryptedPrivateKey)
	if err != nil {
		tx.Rollback()
		return err
	}
	fmt.Printf("private key hash: %x\n", generatedPrivateKeyHash)
	if bytes.Equal(generatedPrivateKeyHash, privateKeyHash) != true {
		tx.Rollback()
		return errors.New("Private keys do not match, wrong password for existing login: " + existingLogin)
	}
	encryptedPrivateKey, err := cryptography.EncryptDataSymmetric([]byte(passwordHashNewLogin), decryptedPrivateKey)
	if err != nil {
		tx.Rollback()
		return err
	}
	u := uuid.New()
	uuidStringBytes, err := u.MarshalText()
	if err != nil {
		tx.Rollback()
		return err
	}
	createTime := time.Now()

	login := Login{Login: newLogin, CreateTime: createTime}
	loginBytes, err := json.Marshal(login)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = loginInformationBucket.Put(uuidStringBytes, loginBytes)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = loginPrivateKeyBucket.Put(uuidStringBytes, encryptedPrivateKey)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func ListAllUsers() {

}
