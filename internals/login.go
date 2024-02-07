package internals

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	bolt "go.etcd.io/bbolt"
	"main/cryptography"
	"time"
)

type Login struct {
	Login      string
	CreateTime time.Time
}

func ListAllLogins() (logins []Login, err error) {
	tx, err := Database.Begin(false)
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}
	loginInfoBucket := tx.Bucket([]byte(BucketLoginInformation))
	if loginInfoBucket == nil {
		return nil, err
	}
	c := loginInfoBucket.Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {
		var login Login
		err = json.Unmarshal(v, &login)
		if err != nil {
			return nil, err
		}
		logins = append(logins, login)

	}
	return logins, nil
}

func CreateFirstLoginWithRSAKeys(username, password string) error {
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
	if err != nil {
		return err
	}
	privateKeyHash, err := cryptography.GetSha256Hash(privateKey)
	if err != nil {
		return err
	}
	privateKeyEncrypted, err := cryptography.EncryptDataSymmetric(passwordHash, privateKey)
	if err != nil {
		return err
	}
	tx, err := Database.Begin(true)
	defer tx.Rollback()
	if err != nil {
		return err
	}
	userInfoBucket := tx.Bucket([]byte(BucketLoginInformation))
	err = userInfoBucket.Put(uuidStringBytes, loginBytes)
	if err != nil {
		return err
	}
	userPrivateKeyEncryptedBucket := tx.Bucket([]byte(BucketLoginPrivateKeyEncrypted))
	err = userPrivateKeyEncryptedBucket.Put(uuidStringBytes, privateKeyEncrypted)
	if err != nil {
		return err
	}
	privateKeyHashBucket := tx.Bucket([]byte(BucketPrivateKeyHash))
	err = privateKeyHashBucket.Put([]byte(PrivateKeyHashKeyName), privateKeyHash)
	if err != nil {
		return err
	}
	publicKeyBucket := tx.Bucket([]byte(BucketPublicKey))
	err = publicKeyBucket.Put([]byte(PublicKeyKeyName), publicKey)
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}
func CreateLoginWithExistingRSAKeys(existingLogin, existingLoginPassword, newLogin, newLoginPassword string) error {
	tx, err := Database.Begin(true)
	defer tx.Rollback()
	if err != nil {
		return err
	}
	passwordHashExistingLogin, err := cryptography.GenerateUserHash([]byte(existingLoginPassword))
	if err != nil {
		return err
	}
	passwordHashNewLogin, err := cryptography.GenerateUserHash([]byte(newLoginPassword))
	if err != nil {
		return err
	}

	userid, err := getLoginId(existingLogin, tx)
	if err != nil {
		return err
	}

	decryptedPrivateKey, err := getAndDecryptPrivateKey(userid, []byte(passwordHashExistingLogin), tx)
	if err != nil {
		return err
	}
	privateKeyHash, err := getPrivateKeyHash(tx)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}
	generatedPrivateKeyHash, err := cryptography.GetSha256Hash(decryptedPrivateKey)
	if err != nil {
		return err
	}
	if bytes.Equal(generatedPrivateKeyHash, privateKeyHash) != true {
		return errors.New("Private keys do not match, wrong password for existing login: " + existingLogin)
	}
	encryptedPrivateKey, err := cryptography.EncryptDataSymmetric([]byte(passwordHashNewLogin), decryptedPrivateKey)
	if err != nil {
		return err
	}
	u := uuid.New()
	uuidStringBytes, err := u.MarshalText()
	if err != nil {
		return err
	}
	createTime := time.Now()

	login := Login{Login: newLogin, CreateTime: createTime}
	loginBytes, err := json.Marshal(login)
	if err != nil {
		return err
	}
	loginInformationBucket := tx.Bucket([]byte(BucketLoginInformation))
	err = loginInformationBucket.Put(uuidStringBytes, loginBytes)
	if err != nil {
		return err
	}
	loginPrivateKeyBucket := tx.Bucket([]byte(BucketLoginPrivateKeyEncrypted))
	err = loginPrivateKeyBucket.Put(uuidStringBytes, encryptedPrivateKey)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
func getLoginId(loginString string, tx *bolt.Tx) (uid []byte, err error) { // return error if no login found
	loginInformationBucket := tx.Bucket([]byte(BucketLoginInformation))
	c := loginInformationBucket.Cursor()

	var userid []byte
	for k, v := c.First(); k != nil; k, v = c.Next() {
		var login Login
		err = json.Unmarshal(v, &login)
		if err != nil {
			return nil, err
		}
		if login.Login == loginString {
			userid = k
		}
	}
	if userid == nil {
		return nil, errors.New("No existing login found: " + loginString)
	}
	return userid, nil
}
